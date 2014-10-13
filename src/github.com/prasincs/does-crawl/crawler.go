package main
import (
    "net/http"
    "code.google.com/p/go.net/html"
    "code.google.com/p/go.net/html/atom"
    "io"
    "fmt"
    "bytes"
    "strings"
    "regexp"
    //"sync"
    //"time"
    //"io/ioutil"
)

// Interfaces influenced heavily by 
type Fetcher interface {
    // Fetch returns the body of URL and
    // a slice of URLs found on that page.
    Fetch(url string) (body string, urls []string, err error)
}

type HttpFetcher struct {

}
 

type Result struct {
    url string
    body string
    links []string
    err error
}
func IsCrawlableUrl(url string, parentUrl string) (bool){
    match, _ := regexp.MatchString("^/", url)
    if (strings.HasPrefix(url, parentUrl) || match){
        return true
    }
    return false
}

    


// Taken and modified from https://gist.github.com/dyoo/6064879 
// The Go playground #69 wasn't quite cutting it
// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
// If the pages start with prefix /, it will try from the parent callee url
func Crawl(url string, depth int, fetcher Fetcher) <-chan Result {
    mapAccess := make(chan map[string]bool, 1)
    mapAccess <- make(map[string]bool)
 
    var loop func(string, int) <-chan Result
    loop = func(url string, depth int) <-chan Result {
        c := make(chan Result)
        go func() {
            if depth <= 0 {
                close(c)
                return
            }
 
            // Check whether we've already visited this.
            m := <-mapAccess
            if _, in := m[url]; in {
                close(c)
                mapAccess <- m
                return
            }
            m[url] = true
            mapAccess <- m
 
            r := fetch(url, fetcher)
            if r.err != nil {
                c <- r
                close(c)
                return
            }

            c <- r
 
            child_chs := make([]<-chan Result, 0)
            for _, u := range r.links {
                if (IsCrawlableUrl(u,url)){
                    u = fmt.Sprintf("%s%s", url, u)
                    child_chs = append(child_chs, loop(u, depth-1))
                }

            }
            for r := range multiplex(child_chs) {
                //match, _ = regexp(MatchString)
                
                c <- r
            }
            close(c)
        }()
        return c
    }
    return loop(url, depth)
}



func fetch(url string, fetcher Fetcher) Result {
    match, _ := regexp.MatchString("^(/|http:)", url)
    if (!match){
        return Result{url, "", nil, nil}
    }
    body, links, err := fetcher.Fetch(url)
    if err != nil {
        return Result{url, "", nil, err}
    } else {
        return Result{url, body, links, nil}
    }
}

func multiplex(chs []<-chan Result) <-chan Result {
    c := make(chan Result)
    d := make(chan bool)
    for _, ch := range chs {
        go func(ch <-chan Result) {
            for r := range ch {
                c <- r
            }
            d <- true
        }(ch)
    }
    go func() {
        for i := 0; i < len(chs); i++ {
            <-d
        }
        close(c)
    }()
    return c
}


func (hf HttpFetcher) Fetch (url string) (string, []string, error){
    res, err := http.Get(url) // Figure out how to DI this
    if err != nil {
        fmt.Printf("ERROR: %s\n", err)
    }
    defer res.Body.Close()
    // TODO - think of less crappy way to do this without copying
    buf := new(bytes.Buffer)
    buf.ReadFrom(res.Body)
    r := strings.NewReader(buf.String())
    urls :=  ExtractUrlsFromHtml(r)
    fmt.Println(urls)
    return buf.String(), urls, nil

}



func ExtractUrlsFromHtml(r io.Reader) ([]string) {
    
    d := html.NewTokenizer(r)
    urls := []string{}
    for { 
        // token type
        tokenType := d.Next() 

        if tokenType == html.ErrorToken {
            return urls
        }       

        token := d.Token()
        switch tokenType {
            case html.StartTagToken: // <tag>
                //fmt.Printf("here: token started %s\n", token.DataAtom)
                if (token.DataAtom == atom.A){
                    for _, a := range token.Attr {
                        if a.Key == "href" {
                            urls = append(urls, a.Val)
                        }
                    }
                }
                // type Token struct {
                //     Type     TokenType
                //     DataAtom atom.Atom
                //     Data     string
                //     Attr     []Attribute
                // }
                //
                // type Attribute struct {
                //     Namespace, Key, Val string
                // }
            case html.TextToken: // text between start and end tag
            case html.EndTagToken: // </tag>
            case html.SelfClosingTagToken: // <tag/>

        }
    }
    return urls
}
