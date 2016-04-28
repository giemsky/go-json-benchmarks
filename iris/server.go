package main

import (
  "strings"
  "log"
  "encoding/json"
  "github.com/kataras/iris"
)

type Record struct {
  Id string
  Timestamp int
  Metadata map[string]interface{}
}

func main() {
  log.Print("HTTP server listening on 8080 port...")

  iris.Post("/", handler)
  iris.Get("/_ah/start", healthHandler)
  iris.Get("/_ah/health", healthHandler)

  iris.Listen()
}

func healthHandler(c *iris.Context) {
  c.Write("OK")
}

func handler(c *iris.Context) {
  var records []Record

  if c.MethodString() != "POST" {
    c.Write("POST method required")
    c.SetStatusCode(400)
    return
  }

  if c.RequestHeader("Content-Type") != "application/json" {
    c.Write("Invalid content type, application/json required")
    c.SetStatusCode(400)
    return
  }

	data := c.RequestCtx.Request.Body()
	decoder := json.NewDecoder(strings.NewReader(string(data)))
	decoder.Decode(&records)

  for _, record := range records {
    json.Marshal(record.Metadata)
    // b, _  := json.Marshal(record.Metadata)
    // log.Print(string(b))
  }

  c.Write("OK")
}
