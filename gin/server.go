package main

import (
  "net/http"
  "log"
  "encoding/json"
  "github.com/gin-gonic/gin"
)

type Record struct {
  Id string
  Timestamp int
  Metadata map[string]interface{}
}

func main() {
  log.Print("HTTP server listening on 8080 port...")
  router := gin.New()

  router.POST("/", handler)
  router.GET("/_ah/start", healthHandler)
  router.GET("/_ah/health", healthHandler)

  router.Run()
}

func healthHandler(c *gin.Context) {
  c.String(http.StatusOK, "OK")
}

func handler(c *gin.Context) {
  var records []Record
  r := c.Request

  if r.Method != "POST" {
    c.String(http.StatusBadRequest, "POST method required")
    return
  }

  if r.Header.Get("Content-Type") != "application/json" {
    c.String(http.StatusBadRequest, "Invalid content type, application/json required")
    return
  }

  decoder := json.NewDecoder(r.Body)
  decoder.Decode(&records)
  for _, record := range records {
    json.Marshal(record.Metadata)
    // b, _  := json.Marshal(record.Metadata)
    // log.Print(string(b))
  }

  c.String(http.StatusOK, "OK")
}
