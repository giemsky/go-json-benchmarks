# Benchmarking JSON handling with Go

## Why?

I was searching for a more efficient server solution as an alternative to [Fluentd](fluentd.org) HTTP endpoint. Since I had previous experience with `Go` language it seemed to be a good candidate for the job. 

The server was supposed to accept `JSON` events in max. 10-records batch. 

Estimated load was around `1500 req/sec`.

## Benchmark

I used [wrk](https://github.com/wg/wrk) tool for benchmark. It seemed more efficient than Apache Benchmark.  

Tests were performed using two linux instances hosted on [DigitalOcean](https://www.digitalocean.com/). 

### Instance parameters
* 2 Core Processor
* 2GB RAM
* 40GB SSD
* Ubuntu 14.04

After a few tests it turned out, that each instance is able to generate 800 req/sec https traffic at maximum. This can be a bottleneck when testing more efficient backends.

Another very important factor is network. Testing from my local machine I got numbers even 6 times less just because of poor coworking network. 

### Command

`wrk -t12 -c50 -d30s -s post.lua https://example.org`

Run on each test instance.

### Curl test POST

`curl -X POST -d post.json -H "Content-Type: application/json" -v https://example.org`

### Hosting

Application was hosten on Google App Engine and Google Compute Engine using `B2` and `n1-highcpu-2` machine types.

## Requirements

* Accept `JSON` request
* Perform `Content-Type` and HTTP Method validation
* Parse `JSON` content
* Encode `Metadata` part back to `JSON`

## Libraries

For HTTP handling I used [fasthttp](https://github.com/valyala/fasthttp), [Gin](https://gin-gonic.github.io/gin/), [Iris](http://kataras.github.io/iris/), [net/http](https://golang.org/pkg/net/http/). `JSON` was handled using [encoding/json](https://golang.org/pkg/encoding/json/) and [ffjson](https://github.com/pquerna/ffjson) libraries.

## Example payload

Payload was generated using online [JSON Generator](http://www.json-generator.com/).

```
[
  {
    "id": "e542e0ec-3b3e-43e3-86ae-c48f57b3fe40",
    "timestamp": 1460616245,
    "metadata": {
      "index": 0,
      "guid": "42b3cde7-c8b5-4739-a277-7b3b0189054f",
      "isActive": true,
      "balance": "$2,784.60",
      "picture": "http://placehold.it/32x32",
      "age": 40,
      "eyeColor": "green",
      "name": "Daniels Parrish",
      "gender": "male",
      "company": "KOZGENE",
      "email": "danielsparrish@kozgene.com",
      "phone": "+1 (968) 562-3911",
      "address": "727 Lorimer Street, Newry, Nevada, 2684",
      "about": "Tempor proident laboris deserunt esse aute occaecat in sit commodo et mollit sunt. Labore laboris dolore nostrud est voluptate tempor fugiat est commodo voluptate esse. Dolore in sit dolor adipisicing ut officia laborum. Ullamco consectetur laboris amet nostrud deserunt qui aliqua cupidatat officia adipisicing ipsum voluptate. Ut Lorem laborum est velit pariatur ut ut velit ipsum cupidatat.\r\n",
      "registered": "2014-01-19T04:20:48 -01:00",
      "latitude": -81.924366,
      "longitude": 80.839932
    }
  },
  {
    "id": "d104c643-9f4f-44a5-9f99-d3603cfc08e7",
    "timestamp": 1460616245,
    "metadata": {
      "index": 1,
      "guid": "9a84f2df-6230-422c-982c-b936392be10f",
      "isActive": false,
      "balance": "$3,216.18",
      "picture": "http://placehold.it/32x32",
      "age": 21,
      "eyeColor": "blue",
      "name": "Katie Harrison",
      "gender": "female",
      "company": "FREAKIN",
      "email": "katieharrison@freakin.com",
      "phone": "+1 (993) 534-3625",
      "address": "198 Jerome Avenue, Thornport, Tennessee, 3677",
      "about": "Enim anim officia id nulla. Ipsum commodo qui minim sit veniam exercitation anim esse ut culpa veniam. Sit elit eiusmod voluptate nisi voluptate eiusmod minim consectetur velit dolor est est cupidatat. Reprehenderit anim sit adipisicing veniam sit veniam anim consectetur est. Amet excepteur id proident excepteur id. Consequat aute eiusmod enim minim. Et ea ipsum eiusmod deserunt.\r\n",
      "registered": "2016-02-16T07:29:30 -01:00",
      "latitude": 3.512001,
      "longitude": -140.639882
    }
  },
  ...
]

```

## Results

I included `Fluentd` just for comparison. It wasn't perfectly accurate, because `Fluentd` had more complicated implementation (e.g. pushed parsed `JSON` events to BigQuery). I just wanted to see how much `Go` version is faster and then implement other missing features for stage 2 testing if results were promising.

| Tool | Result |
|---|---|
| fasthttp | 907 req/sec |
| GAE | 407 req/sec |
| Gin | 790 req/sec |
| Iris | 795 req/sec |
| net/http + ffjson | 803 req/sec |
| net/http + encoding/json | 836 req/sec |
| Fluentd | 572 req/sec |

The following results were generated from two tests running simultaneously for 30s.

### fasthttp

```
  12 threads and 50 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   108.42ms   50.11ms   1.35s    94.92%
    Req/Sec    37.83      9.85    70.00     72.86%
  13522 requests in 30.03s, 3.34MB read
Requests/sec:    450.23
Transfer/sec:    113.88KB
```

```
  12 threads and 50 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   104.83ms   28.65ms 328.67ms   81.89%
    Req/Sec    38.45      9.64    80.00     73.05%
  13742 requests in 30.05s, 3.35MB read
Requests/sec:    457.33
Transfer/sec:    114.33KB
```

### Gin

```
  12 threads and 50 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   124.59ms   53.33ms 702.30ms   91.59%
    Req/Sec    33.71      9.91    80.00     77.33%
  11835 requests in 30.05s, 2.90MB read
Requests/sec:    393.79
Transfer/sec:     98.83KB
```

```
  12 threads and 50 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   123.92ms   53.64ms 704.56ms   91.70%
    Req/Sec    34.02      9.66    80.00     79.74%
  11904 requests in 30.06s, 2.92MB read
Requests/sec:    396.03
Transfer/sec:     99.39KB
```

### Iris

```
  12 threads and 50 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   120.92ms   34.87ms 564.89ms   83.25%
    Req/Sec    33.45      9.64    80.00     77.99%
  11928 requests in 30.05s, 2.92MB read
Requests/sec:    396.98
Transfer/sec:     99.63KB
```

```
Running 30s test @ https://go-iris-dot-mychat-tracking-beta-1245.appspot.com
  12 threads and 50 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   120.72ms   42.92ms   1.58s    88.30%
    Req/Sec    33.57      9.26    80.00     79.17%
  11977 requests in 30.04s, 2.94MB read
Requests/sec:    398.65
Transfer/sec:    100.05KB
```

### net/http + ffjson

```
  12 threads and 50 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   120.51ms   42.41ms   1.30s    88.59%
    Req/Sec    33.95      9.46    80.00     79.26%
  12020 requests in 30.04s, 2.96MB read
Requests/sec:    400.11
Transfer/sec:    100.81KB
```

```
  12 threads and 50 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   119.06ms   34.02ms 351.91ms   84.15%
    Req/Sec    34.17      9.52    70.00     78.27%
  12111 requests in 30.08s, 2.98MB read
Requests/sec:    402.58
Transfer/sec:    101.43KB
```


### net/http + encoding/json

```
  12 threads and 50 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   126.70ms  132.14ms   1.91s    95.47%
    Req/Sec    35.62     10.31    80.00     75.00%
  12553 requests in 30.05s, 3.10MB read
  Socket errors: connect 0, read 0, write 0, timeout 19
Requests/sec:    417.69
Transfer/sec:    105.65KB
```

```
  12 threads and 50 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   128.77ms  141.78ms   1.96s    95.58%
    Req/Sec    35.72     10.39    70.00     73.58%
  12571 requests in 30.04s, 3.11MB read
  Socket errors: connect 0, read 0, write 0, timeout 12
Requests/sec:    418.42
Transfer/sec:    105.83KB
```

### Google App Engine

```
  12 threads and 50 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    91.03ms   99.14ms   1.55s    97.83%
    Req/Sec    21.93     12.24    80.00     80.98%
  6369 requests in 30.05s, 1.45MB read
  Socket errors: connect 0, read 0, write 0, timeout 193
Requests/sec:    211.91
Transfer/sec:     49.25KB
```

```
  12 threads and 50 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    86.32ms   67.16ms   1.50s    98.17%
    Req/Sec    20.99     11.68    80.00     84.55%
  5858 requests in 30.05s, 1.33MB read
  Socket errors: connect 0, read 0, write 0, timeout 197
Requests/sec:    194.97
Transfer/sec:     45.32KB
```

### Fluentd

```
  12 threads and 50 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   168.23ms   48.23ms   1.19s    82.59%
    Req/Sec    24.37      9.97    60.00     65.21%
  8584 requests in 30.06s, 5.50MB read
Requests/sec:    285.58
Transfer/sec:    187.41KB
```

```
  12 threads and 50 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   166.72ms   44.58ms 488.04ms   79.04%
    Req/Sec    24.37      9.54    70.00     69.40%
  8626 requests in 30.05s, 5.53MB read
Requests/sec:    287.02
Transfer/sec:    188.36KB
```

## Summary

I was disappointed by the test results. `JSON` handling seems to be very slow. It seems to be the weakest `Go` feature, so often used in REST world. There is an [open ticket](https://github.com/golang/go/issues/5683) for this matter. 

### Alternatives

If you don't need `JSON` for communication you can consider the following alternatives.

* [Protocol Buffers](https://developers.google.com/protocol-buffers/)
* [gob](https://golang.org/pkg/encoding/gob/)

These should be blazing fast according to what I've read.