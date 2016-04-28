package main

import (
  "fmt"
  "net/http"
  "encoding/json"
)

type Record struct {
  Id string
  Timestamp int
  Metadata map[string]interface{}
}

func init() {
  http.HandleFunc("/", handler)
  http.HandleFunc("/_ah/start", healthHandler)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprint(w, "OK")
}

func handler(w http.ResponseWriter, r *http.Request) {
  var records []Record

  if r.Method != "POST" {
    w.WriteHeader(http.StatusBadRequest)
    fmt.Fprint(w, "POST method required")
    return
  }

  if r.Header.Get("Content-Type") != "application/json" {
    w.WriteHeader(http.StatusBadRequest)
    fmt.Fprint(w, "Invalid content type, application/json required")
    return
  }

  decoder := json.NewDecoder(r.Body)
  decoder.Decode(&records)
  for _, record := range records {
    json.Marshal(record.Metadata)
  }

  fmt.Fprint(w, "OK")
}
