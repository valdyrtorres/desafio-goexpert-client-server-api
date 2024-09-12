package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	req, err := http.NewRequest(http.MethodGet, "http://httpbin.org/get", nil)
	if err != nil {
		log.Fatal(err)
	}

	// 80ms
	// we have an HTTP request and we want to control the timeout of such request.
	// Ok, what if we want to make sure this request completes within 80ms?
	// We just need to use context.WithTimeout!
	// if we put 80ms -> 2024/09/12 10:12:59 Get "http://httpbin.org/get": context deadline exceeded
	// if we put 800ms -> the result is properly
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Millisecond*800))
	defer cancel()
	req = req.WithContext(ctx)

	c := &http.Client{}
	res, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	out, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(out))
}
