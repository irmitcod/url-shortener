package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"time"
)

var (
	reqs int
	max  int
)

func init() {
	flag.IntVar(&reqs, "reqs", 10000, "Total requests")
	flag.IntVar(&max, "concurrent", 100, "Maximum concurrent requests")
}

type Response struct {
	*http.Response
	err error
}

// Dispatcher
func dispatcher(reqChan chan *http.Request) {
	defer close(reqChan)
	for i := 0; i < reqs; i++ {

		form := url.Values{}
		form.Add("OriginalURL", fmt.Sprintf("https://google.com/?search=urls%d", i))
		req, err := http.NewRequest("POST", "http://localhost:8080/UrlShotener", strings.NewReader(form.Encode()))

		if err != nil {
			log.Println(err)
		}

		req.Header = http.Header{
			"Content-Type":  []string{"application/x-www-form-urlencoded"},
			"Authorization": []string{"Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiNjI3Y2FlYjI5MWM1NjRjMjg0YmUxYzNjIiwiaXNfYWRtaW4iOmZhbHNlLCJleHAiOjE2NTIzNDQzNjksImp0aSI6IjYyN2NhZWIyOTFjNTY0YzI4NGJlMWMzYyJ9.jRHeyOR2uuAVUXarIC-X_-vvyLbQOQ2ggK2QcJBG73hFOy8mAzk6o-LYmm_NO7OqGOys6Z_Geryrjgw6AQL7Gw"},
		}
		reqChan <- req
	}
}

// Worker Pool
func workerPool(reqChan chan *http.Request, respChan chan Response) {
	t := &http.Transport{}
	for i := 0; i < max; i++ {
		go worker(t, reqChan, respChan)
	}
}

// Worker
func worker(t *http.Transport, reqChan chan *http.Request, respChan chan Response) {
	for req := range reqChan {
		resp, err := t.RoundTrip(req)
		r := Response{resp, err}
		respChan <- r
	}
}

// Consumer
func consumer(respChan chan Response) (int64, int64) {
	var (
		conns int64
		size  int64
	)
	for conns < int64(reqs) {
		select {
		case r, ok := <-respChan:

			if ok {
				if r.err != nil {
					log.Println(r.err)
				} else {
					size += r.ContentLength
					//log.Println(r.Response.Body)
					if err := r.Body.Close(); err != nil {
						log.Println(r.err)
					}
				}
				conns++
			}
		}
	}
	return conns, size
}

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	//c := config.GetConfig()
	//
	//_,keys := config.NewMemoryClient(c)

	//log.Println(keys)

	reqChan := make(chan *http.Request)
	respChan := make(chan Response)
	start := time.Now()
	go dispatcher(reqChan)
	go workerPool(reqChan, respChan)
	conns, size := consumer(respChan)
	took := time.Since(start)
	ns := took.Nanoseconds()
	av := ns / conns
	average, err := time.ParseDuration(fmt.Sprintf("%d", av) + "ns")
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("Connections:\t%d\nConcurrent:\t%d\nTotal size:\t%d bytes\nTotal time:\t%s\nAverage time:\t%s\n", conns, max, size, took, average)
}
