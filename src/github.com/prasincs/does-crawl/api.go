package main

import (
	"fmt"
	"github.com/codegangsta/martini"
	"net/http"
	"strconv"
	"time"
	// "reflect"
)

func GetUrls(r *http.Request, enc Encoder, db DB) string {
	// Get the query string arguments, if any
	qs := r.URL.Query()
	link, parent := qs.Get("link"), qs.Get("parent")
	if link != "" || parent != "" {
		// At least one filter, use Find()
		return Must(enc.Encode(toIface(db.Find(link, parent))...))
	}
	// Otherwise, return all albums
	return Must(enc.Encode(toIface(db.GetAll())...))
}

func GetUrl(enc Encoder, db DB, parms martini.Params) (int, string) {
	id, err := strconv.Atoi(parms["id"])
	al := db.Get(id)
	if err != nil || al == nil {
		// Invalid id, or does not exist
		return http.StatusNotFound, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("the url with id %s does not exist", parms["id"]))))
	}
	return http.StatusOK, Must(enc.Encode(al))
}

// UpdateAlbum changes the specified album.
// The order is messed up, it does request before inserting or checking - FIX
func UpdateUrl(r *http.Request, enc Encoder, db DB, parms martini.Params) (int, string) {
	urlParams, err := getPutPostParams(r)
	if err != nil {
		// Invalid id, 404
		return http.StatusNotFound, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("the url with id %s does not exist", urlParams.Id))))
	}
	u, err2 := doCrawl(urlParams.Link, urlParams.Parent, urlParams.MaxDepth)
	if err2 != nil {
		return http.StatusInternalServerError, Must(enc.Encode(NewError(ErrCrawlingFailed, fmt.Sprintf("the url %s failed to be crawled. Error: %s", err2))))
	}
	err = db.Update(u)
	switch err {
	case ErrAlreadyExists:
		return http.StatusConflict, Must(enc.Encode(
			NewError(ErrCodeAlreadyExists, fmt.Sprintf("the url '%s' from '%s' already exists", u.Link, u.Parent))))
	case nil:
		return http.StatusOK, Must(enc.Encode(u))
	default:
		panic(err)
	}
}

func AddUrl(w http.ResponseWriter, r *http.Request, enc Encoder, db DB) (int, string) {
	urlParams, err := getPutPostParams(r)
	if err != nil {
		return http.StatusBadRequest, Must(enc.Encode(
			NewError(ErrInvalidValueSupplied, fmt.Sprintf("the input supplied was unacceptable: %s", err))))
	}
	u, err2 := doCrawl(urlParams.Link, urlParams.Parent, urlParams.MaxDepth)
	if err2 != nil {
		return http.StatusInternalServerError, Must(enc.Encode(NewError(ErrCrawlingFailed, fmt.Sprintf("the url %s failed to be crawled. Error: %s", err2))))
	}
	id, err := db.Add(u)
	switch err {
	case ErrAlreadyExists:
		// Duplicate
		return http.StatusConflict, Must(enc.Encode(
			NewError(ErrCodeAlreadyExists, fmt.Sprintf("the url '%s' from '%s' already exists", u.Link, u.Parent))))
	case nil:
		// TODO : Location is expected to be an absolute URI, as per the RFC2616
		w.Header().Set("Location", fmt.Sprintf("/urls/%d", id))
		return http.StatusCreated, Must(enc.Encode(u))
	default:
		panic(err)
	}
}

var httpFetcher = HttpFetcher{}

func getPutPostParams(r *http.Request) (*Url, error) {
	link, parent, maxDepthStr, idStr := r.FormValue("link"), r.FormValue("parent"), r.FormValue("maxDepth"), r.FormValue("id")
	// Note: this is a terrible place to do the blocking crawling..
	id := 0
	if maxDepthStr == "" {
		maxDepthStr = "4"
	}
	maxDepth, err := strconv.Atoi(maxDepthStr)
	if err != nil {
		return nil, err
	}

	if idStr != "" {
		id, err = strconv.Atoi(idStr)
		if err != nil {
			return nil, err
		}
	}

	return &Url{
		Id:       id,
		Link:     link,
		Parent:   parent,
		MaxDepth: maxDepth,
	}, nil
}

func doCrawl(link string, parent string, maxDepth int) (*Url, error) {
	links := []CrawledLink{}
	crawler := Crawl(link, maxDepth, httpFetcher)
	for req := range crawler {
		if req.err != nil {
			fmt.Println(req.err)
		} else {
			for _, childLink := range req.links {
				// fmt.Println(reflect.TypeOf(childLink))
				crawledLink := &CrawledLink{
					Link:   childLink,
					Parent: req.url,
					Depth:  (maxDepth - req.depth),
				}
				links = append(links, *crawledLink)
			}
			// fmt.Println(reflect.TypeOf(req.links))
			//fmt.Printf("found: %s %q\n", req.url, req.links)
		}
	}
	// assumes that it's synchronous, but I think it's sufficient for now
	t := time.Now().Local()
	return &Url{
		Link:        link,
		Parent:      parent,
		Links:       links,
		MaxDepth:    maxDepth,
		LastCrawled: t.Format(time.RFC3339),
	}, nil
}

func toIface(v []*Url) []interface{} {
	if len(v) == 0 {
		return nil
	}
	ifs := make([]interface{}, len(v))
	for i, v := range v {
		ifs[i] = v
	}
	return ifs
}
