package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

func main() {
	server := []string{
		"http://localhost:8083",
		"http://localhost:8081",
		"http://localhost:8082",
	}
	for {
		before := time.Now()
		//res := Get(server[0])
		//res := Read(server[0], time.Second)
		res := MultiRead(server, time.Second)
		after := time.Now()
		fmt.Println("Response:", *res)
		fmt.Println("Time:", after.Sub(before))
		fmt.Println()
		time.Sleep(500 * time.Millisecond)
	}
}

type Response struct {
	Body       string
	StatusCode int
}

// Get makes an HTTP Get request and returns an abbreviated response.
// Status code 200 means that the request was successful.
// The function returns &Response{"", 0} if the request fails
// and it blocks forever if the server doesn't respond.
func Get(url string) *Response {
	res, err := http.Get(url)
	if err != nil {
		return &Response{}
	}
	// res.Body != nil when err == nil
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("ReadAll: %v", err)
	}
	return &Response{string(body), res.StatusCode}
}

// FIXME
// I've found two insidious bugs in this function; both of them are unlikely
// to show up in testing. Please fix them right away – and don't forget to
// write a doc comment this time.
func Read(url string, timeout time.Duration) *Response {
	response := make(chan *Response)
	go func() {
		response <- Get(url)
	}()
	select {
	case res := <-response:
		return res
	case <-time.After(timeout):
		return &Response{"Gateway timeout\n", 504}
	}
	return nil
}

// MultiRead makes an HTTP Get request to each url and returns
// the response of the first server to answer with status code 200.
// If none of the servers answer before timeout, the response is
// 503 – Service unavailable.
func MultiRead(urls []string, timeout time.Duration) (res *Response) {
	response := make(chan *Response, len(urls))
	wg := new(sync.WaitGroup)
	wg.Add(len(urls))
	for _, url := range urls {
		go func(url string) {
			read := Read(url, timeout)
			if read.StatusCode == 200 {
				response <- read
			}
			wg.Done()
			wg.Wait()
			response <- &Response{"Service Unavailable", 503}
		}(url)
	}
	res = <-response
	return
}
