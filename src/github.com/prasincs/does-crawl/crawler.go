package main
import (
    "net/http"
    "code.google.com/p/go.net/html"
    "code.google.com/p/go.net/html/atom"
    "strings"
    "io"
    "fmt"
    "sync"
    //"time"
    //"io/ioutil"
)
type Visited struct {
    url_map map[Url]bool
    mutex sync.Mutex
}

func (self *Visited) TestAndSetVisited(url string) bool {
    defer func() {
        self.url_map[url] = true
        self.mutex.Unlock()
    }()
    
    self.mutex.Lock()
    return self.url_map[url]
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher, 
           visited *Visited, wg *sync.WaitGroup) {
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
                child_chs = append(child_chs, loop(u, depth-1))
            }
            for r := range multiplex(child_chs) {
                c <- r
            }
            close(c)
        }()
        return c
    }
    return loop(url, depth)
}

type Result struct {
    url string
    body string
    links []string
    err error
}


func (f Fetcher) Fetch(url string) (string, []string,
                                        error) {
    if res, ok := f[url]; ok {
        return res.body, res.urls, nil
    }
    return "", nil, fmt.Errorf("not found: %s", url)
}




func extractUrlsFromHtml(r io.Reader, parentUrl string, foundUrls chan string) {
    
    d := html.NewTokenizer(r)
    for { 
        // token type
        tokenType := d.Next() 
        if tokenType == html.ErrorToken {
            return     
        }       
        token := d.Token()
        //fmt.Println(token);
        switch tokenType {
            case html.StartTagToken: // <tag>
                if (token.DataAtom == atom.A){
                    for _, a := range token.Attr {
                        if a.Key == "href" {
                            fmt.Printf("Found %s, Sending via channel\n",a.Val)
                            foundUrls <- a.Val
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
}

func processFoundUrls(urlsQueryChan chan string, foundUrls chan string){
    for{
           urlFound := <- foundUrls
           fmt.Printf("Pocessing: %s\n", urlFound)
           if (strings.HasPrefix(urlFound,"/")){
                urlsQueryChan <- urlFound
           }
    }
}

func httpQueryToBody(urlsQueryChan chan string, foundUrls chan string){
    fmt.Println("here");
    for {
        parentUrl := <-urlsQueryChan
        resp, err := http.Get(<-urlsQueryChan)
        if err != nil {
            fmt.Printf("ERROR: %s\n", err)
        }
        defer resp.Body.Close()
        extractUrlsFromHtml(resp.Body, parentUrl, foundUrls)
        }
}

