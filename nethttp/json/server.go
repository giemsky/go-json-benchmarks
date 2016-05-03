package main

import (
  "fmt"
  "net/http"
  "log"
  "encoding/json"
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

  decoder := json.NewDecoder(r.Body)
  decoder.Decode(&records)
  for _, record := range records {
    json.Marshal(record.Metadata)
    // b, _  := json.Marshal(record.Metadata)
    // log.Print(string(b))
  }

  fmt.Fprint(w, "OK")
}
