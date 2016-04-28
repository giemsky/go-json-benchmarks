package main

import (
  "fmt"
  "github.com/valyala/fasthttp"
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
  if err := fasthttp.ListenAndServe(":8080", requestHandler); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}

func requestHandler(ctx *fasthttp.RequestCtx) {
  var records []Record

  if string(ctx.Method()) != "POST" {
    ctx.Response.SetStatusCode(400)
    fmt.Fprint(ctx, "POST method required")
    return
  }

  if string(ctx.Request.Header.ContentType()) != "application/json" {
    ctx.Response.SetStatusCode(400)
    fmt.Fprint(ctx, "Invalid content type, application/json required")
    return
  }

  json.Unmarshal(ctx.PostBody(), &records)
  for _, record := range records {
    json.Marshal(record.Metadata)
    // b, _ := json.Marshal(record.Metadata)
    // log.Print(string(b))
  }

  fmt.Fprint(ctx, "OK")
}
