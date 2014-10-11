package main
import (
    "net/http"
    "code.google.com/p/go.net/html"
    "code.google.com/p/go.net/html/atom"
    "strings"
    "io"
    "fmt"
    "time"
    //"io/ioutil"
)

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

