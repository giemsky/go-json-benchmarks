package main

import (
  "fmt"
  "net/http"
  "log"
  "github.com/pquerna/ffjson/ffjson"
)

type Record struct {
  Id string
  Timestamp int
  Metadata map[string]interface{}
}

func main() {
  log.Print("HTTP server listening on 8080 port...")
  http.HandleFunc("/", handler)
  http.HandleFunc("/_ah/start", healthHandler)
  http.HandleFunc("/_ah/health", healthHandler)
  http.ListenAndServe(":8080", nil)
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

  decoder := ffjson.NewDecoder()
  decoder.DecodeReader(r.Body, &records)
  for _, record := range records {
    ffjson.Marshal(record.Metadata)
    // b, _  := ffjson.Marshal(record.Metadata)
    // log.Print(string(b))
  }

  fmt.Fprint(w, "OK")
}
