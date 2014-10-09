package main

import (
  "fmt"
  "net/http"
  "strconv"

  "github.com/codegangsta/martini"
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
      NewError(ErrCodeNotExist, fmt.Sprintf("the album with id %s does not exist", parms["id"]))))
  }
  return http.StatusOK, Must(enc.Encode(al))
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

