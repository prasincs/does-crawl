package main

import (
	//"encoding/xml"
	"errors"
	//"fmt"
	"sync"
)

// The DB interface defines methods to manipulate the urls
type DB interface {
	Get(id int) *Url
	GetAll() []*Url
	Find(link, parent string) []*Url
	Add(a *Url) (int, error)
	Update(a *Url) error
	Delete(id int)
}

var (
	ErrAlreadyExists = errors.New("album already exists")
)

// The one and only database instance.
var db DB

func init() {
	db = &urlsDB{
		m: make(map[int]*Url),
	}
	// Fill the database

	// TODO move to tests
	//db.Add(&Url{Id: 1, Link: "google.com", Parent: ""})
	//db.Add(&Url{Id: 2, Link: "google.com/mail", Parent: "google.com"})
}

type CrawledLink struct {
	Link   string `json: "link"`
	Parent string `json: "parent"`
	Depth  int    `json: "depth"`
}

type Url struct {
	Id          int           `json: "id"`
	Link        string        `json: "link"`
	Parent      string        `json: "parent"`
	MaxDepth    int           `json: "maxDepth"`
	Links       []CrawledLink `json: "links"`
	LastCrawled string        `json: "lastCrawled"`
}

// Thread-safe in-memory map of urls.
type urlsDB struct {
	sync.RWMutex
	m   map[int]*Url
	seq int
}

// GetAll returns all urls from the database.
func (db *urlsDB) GetAll() []*Url {
	db.RLock()
	defer db.RUnlock()
	if len(db.m) == 0 {
		return nil
	}
	ar := make([]*Url, len(db.m))
	i := 0
	for _, v := range db.m {
		ar[i] = v
		i++
	}
	return ar
}

// Get returns the url identified by the id, or nil.
func (db *urlsDB) Get(id int) *Url {
	db.RLock()
	defer db.RUnlock()
	return db.m[id]
}

// Add creates a new url and returns its id, or an error.
func (db *urlsDB) Add(a *Url) (int, error) {
	db.Lock()
	defer db.Unlock()
	// Return an error if band-title already exists
	if !db.isUnique(a) {
		return 0, ErrAlreadyExists
	}
	// Get the unique ID
	db.seq++
	a.Id = db.seq
	// Store
	db.m[a.Id] = a
	return a.Id, nil
}

// Update changes the url identified by the id. It returns an error if the
// updated url is a duplicate.
func (db *urlsDB) Update(a *Url) error {
	db.Lock()
	defer db.Unlock()
	if !db.isUnique(a) {
		return ErrAlreadyExists
	}
	db.m[a.Id] = a
	return nil
}

// Delete removes the url identified by the id from the database. It is a no-op
// if the id does not exist.
func (db *urlsDB) Delete(id int) {
	db.Lock()
	defer db.Unlock()
	delete(db.m, id)
}

// Find returns albums that match the search criteria.
func (db *urlsDB) Find(link, parent string) []*Url {
	db.RLock()
	defer db.RUnlock()
	var res []*Url
	for _, v := range db.m {
		if v.Link == link || link == "" {
			if v.Parent == parent || parent == "" {
				res = append(res, v)
			}
		}
	}
	return res
}

func (db *urlsDB) isUnique(a *Url) bool {
	for _, v := range db.m {
		if v.Link == a.Link && v.Parent == a.Parent && v.Id != a.Id {
			return false
		}
	}
	return true
}
