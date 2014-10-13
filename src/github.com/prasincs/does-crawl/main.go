package main

import (
  "net/http"
  "strings"
  "regexp"
  "github.com/codegangsta/martini"
  "fmt"
)

const (
  numCrawlers = 2
)

var urls = []string{
  "http://www.google.com/",
  "http://golang.org/",
  "http://blog.golang.org/",
}



var m *martini.Martini

func init() {
  m = martini.New()
  // Setup middleware
  m.Use(martini.Recovery())
  m.Use(martini.Logger())

  m.Use(MapEncoder)
  // Setup routes
  r := martini.NewRouter()
  r.Get(`/urls`, GetUrls)
  r.Get(`/urls/:id`, GetUrl)
  r.Post(`/urls`, AddUrl)
  //r.Put(`/urls/:id`, UpdateUrl)
  //r.Delete(`/urls/:id`, DeleteUrl)
  // Inject database
  m.MapTo(db, (*DB)(nil))
  // Add the router action
  m.Action(r.Handle)
}

//Taken from https://github.com/PuerkitoBio/martini-api-example/blob/master/server.go

// The regex to check for the requested format (allows an optional trailing
// slash).
var rxExt = regexp.MustCompile(`(\.(?:xml|text|json))\/?$`)

// MapEncoder intercepts the request's URL, detects the requested format,
// and injects the correct encoder dependency for this request. It rewrites
// the URL to remove the format extension, so that routes can be defined
// without it.
func MapEncoder(c martini.Context, w http.ResponseWriter, r *http.Request) {
  // Get the format extension
  matches := rxExt.FindStringSubmatch(r.URL.Path)
  ft := ".json"
  if len(matches) > 1 {
    // Rewrite the URL without the format extension
    l := len(r.URL.Path) - len(matches[1])
    if strings.HasSuffix(r.URL.Path, "/") {
      l--
    }
    r.URL.Path = r.URL.Path[:l]
    ft = matches[1]
  }
  // Inject the requested encoder
  switch ft {
  case ".xml":
    c.MapTo(xmlEncoder{}, (*Encoder)(nil))
    w.Header().Set("Content-Type", "application/xml")
  case ".text":
    c.MapTo(textEncoder{}, (*Encoder)(nil))
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
  default:
    c.MapTo(jsonEncoder{}, (*Encoder)(nil))
    w.Header().Set("Content-Type", "application/json")
  }
}

func main() {
  httpFetcher := HttpFetcher{}
  c := Crawl("http://google.com", 4, httpFetcher)
  for r := range c {
    if r.err != nil {
      fmt.Println(r.err)
    } else {
      fmt.Printf("found: %s %q\n", r.url, r.links)
    }
  }
  // s := "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Transitional//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd\">\n<html xmlns=\"http://www.w3.org/1999/xhtml\">\n<head>\n<meta http-equiv=\"Content-Type\" content=\"text/html; charset=iso-8859-1\" />\n<title>HostGator Web Hosting Website Startup Guide</title>\n<link rel=\"stylesheet\" href=\"http://www.hostgatorsupport.com/style.css\" type=\"text/css\" />\n</head>\n\n<body>\n<br />\n<div id=\"wrap\">\n\n<div id=\"header\"><img src=\"http://www.hostgatorsupport.com/images/ban2.png\" alt=\"ban\" width=\"768\" height=\"141\" /></div>\n\n\n<div id=\"content\"><a href=\"/webmail\" target=\"_blank\"></a><br />\n\n  <br />\n  <table width=\"100%\" border=\"0\">\n  <tr>\n    <td align=\"right\"><img src=\"http://www.hostgatorsupport.com/images/cp.png\" alt=\"cp\" width=\"51\" height=\"40\" /></td>\n    <td valign=\"middle\" class=\"contentlinks\"><a href=\"/cpanel\" target=\"_blank\">cPanel Login</a></td>\n    <td align=\"right\"><img src=\"http://www.hostgatorsupport.com/images/mail2.png\" alt=\"webmail\" width=\"51\" height=\"36\" /></td>\n    <td class=\"contentlinks\"><a href=\"/webmail\" target=\"_blank\">Webmail Login</a></td>\n  </tr>\n\n  <tr>\n    <td width=\"24%\" align=\"right\"><img src=\"http://www.hostgatorsupport.com/images/start.png\" alt=\"start\" width=\"45\" height=\"35\" /></td>\n        <td width=\"24%\" valign=\"middle\"><a href=\"http://hostgator.com/gettingstarted.shtml\" target=\"_blank\" class=\"contentlinks\">Getting Started</a></td>\n        <td width=\"8%\" align=\"right\"><img src=\"http://www.hostgatorsupport.com/images/dollar.png\" alt=\"dollar\" width=\"35\" height=\"40\" /></td>\n        <td width=\"40%\"><a href=\"https://secure.hostgator.com/billing\" target=\"_blank\" class=\"contentlinks\">Billing / Invoices </a></td>\n      </tr>\n  <tr>\n    <td align=\"right\"><img src=\"http://www.hostgatorsupport.com/images/filmstrip2.jpg\" alt=\"film2\" width=\"46\" height=\"39\" />&nbsp;</td>\n\n        <td><a href=\"http://www.hostgator.com/tutorials.shtml\" target=\"_blank\" class=\"contentlinks\">Video Tutorials</a></td>\n        <td align=\"right\"><img src=\"http://www.hostgatorsupport.com/images/com.png\" alt=\"com\" width=\"45\" height=\"35\" /></td>\n        <td><a href=\"http://www.hostgator.com/domains\" target=\"_blank\" class=\"contentlinks\">Purchase / Transfer Domain Name </a></td>\n      </tr>\n  <tr>\n    <td align=\"right\"><img src=\"http://www.hostgatorsupport.com/images/book3.png\" alt=\"book\" width=\"42\" height=\"42\" />&nbsp;&nbsp;</td>\n        <td class=\"contentlinks\"><a href=\"http://support.hostgator.com/index.php?_m=knowledgebase&amp;_a=view\" target=\"_blank\">Knowledgebase</a></td>\n\n        <td align=\"right\"><img src=\"http://www.hostgatorsupport.com/images/mail1.png\" alt=\"mail\" width=\"45\" height=\"48\" /></td>\n        <td><a href=\"http://support.hostgator.com/\" target=\"_blank\" class=\"contentlinks\">Ticket System </a></td>\n      </tr>\n  <tr>\n    <td align=\"right\"><img src=\"http://www.hostgatorsupport.com/images/people.png\" alt=\"people\" width=\"51\" height=\"40\" /></td>\n        <td class=\"contentlinks\"><a href=\"http://forums.hostgator.com\" target=\"_blank\">Online Forums</a></td>\n        <td align=\"right\"><img src=\"http://www.hostgatorsupport.com/images/phone2.png\" alt=\"phone2\" width=\"47\" height=\"45\" /></td>\n        <td class=\"contentlinks\"><a href=\"http://www.hostgator.com/contact.shtml\" target=\"_blank\">Contact US</a></td>\n\n      </tr>\n  </table>\n    <table width=\"100%\" border=\"0\">\n      <tr>\n        <td align=\"left\">&nbsp;</td>\n      </tr>\n    </table>\n</div>\n<div id=\"footer\"><a href=\"javascript:void(0);\" onclick=\"window.open('http://chat.hostgator.com/liveperson/', (Math.floor(Math.random()*100000)), 'toolbar=0,scrollbars=1,location=0,statusbar=0,menubar=0,resizable=0,width=500,height=375,left = 310,top = 275');\"><img src=\"http://www.hostgatorsupport.com/images/banner1.jpg\" alt=\"Live Support\" border=\"0\" /></a><br />\n    | <a href=\"http://www.hostgator.com/\" title=\"HostGator Web Hosting\" target=\"_blank\">HostGator.com Web Hosting</a> |<br />\n\n  <span class=\"style8\">Copyright 2009 &copy; HostGator.com</span></div>\n\n</div>\n<script language='javascript' type=\"JavaScript/text\" src='http://chat.hostgator.com/liveperson/'> </script>\n</body>\n</html>\n\n"
  // r := strings.NewReader(s)
  // urls = ExtractUrlsFromHtml(r)
  // fmt.Println(urls)
  //http.ListenAndServe(":9999", m)
}
