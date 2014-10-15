package main

import (
	"flag"
	"fmt"
	"github.com/codegangsta/martini"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var m *martini.Martini

func init() {
	log.Println("Initializing configuration")
	configuration := GetConfiguration("conf/conf.json")
	m = martini.New()
	// Setup middleware
	m.Use(martini.Recovery())
	m.Use(martini.Logger())
	static := martini.Static("public", martini.StaticOptions{Fallback: "/index.html", Exclude: "/api/v"})
	m.Use(static)
	m.Use(MapEncoder)
	// Setup routes
	r := martini.NewRouter()
	r.Get(`/api/v1/urls`, GetUrls)
	r.Get(`/api/v1/urls/:id`, GetUrl)
	r.Post(`/api/v1/urls`, AddUrl)
	//r.Put(`/urls/:id`, UpdateUrl)
	//r.Delete(`/urls/:id`, DeleteUrl)
	if configuration.Db.Type == "memory" {
		// Inject database
		m.MapTo(db, (*DB)(nil))

	} else {
		log.Println("Database type not configured.")
	}
	// TODO - add other types of databases
	// Add the router action
	m.Action(r.Handle)
}

//Taken from https://github.com/PuerkitoBio/martini-api-example/blob/master/server.go
// The regex to check for the requested format (allows an optional trailing
// slash).
var rxExt = regexp.MustCompile(`(\.(?:text|json))\/?$`)

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
	case ".text":
		c.MapTo(textEncoder{}, (*Encoder)(nil))
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	default:
		c.MapTo(jsonEncoder{}, (*Encoder)(nil))
		w.Header().Set("Content-Type", "application/json")
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [inputfile]\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {

	portPtr := flag.Int("port", 9999, "an int")
	flag.Parse()
	fmt.Printf("Starting server at port %d\n", *portPtr)

	err := http.ListenAndServe(fmt.Sprintf(":%d", *portPtr), m)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
