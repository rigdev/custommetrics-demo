package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func getRequestTime() time.Duration {
	requestTime := time.Second
	s := os.Getenv("REQUEST_TIME")
	if s != "" {
		d, err := time.ParseDuration(s)
		if err != nil {
			fmt.Println("Error parsing request duruation", err.Error())
		} else {
			requestTime = d
		}
	}
	return requestTime
}

func main() {
	for {
		// A hack to not reuse the same connection for each request
		// Otherwise the same consumer gets hit over and over
		client := http.Client{}
		http.DefaultTransport.(*http.Transport).DisableKeepAlives = true

		requestURL := "http://consumer:2112/consume"
		if _, err := client.Get(requestURL); err != nil {
			fmt.Printf("error making http request: %s\n", err)
		}

		time.Sleep(getRequestTime())
	}
}
