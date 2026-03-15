package main

import (
	"fmt"
	"net/http"
)

type Result struct {
	Response *http.Response
	Err      error
}

func errorHandlingExample() {
	maxErrs := 3
	errCount := 0
	done := make(chan any)
	defer close(done)
	for res := range checkStatus(done, "a", "b", "c", "d") {
		if res.Err == nil {
			fmt.Printf("response status: %s\n", res.Response.Status)
			continue
		}
		errCount++
		if errCount == maxErrs {
			fmt.Printf("max errors reached, exiting...\n")
			break
		}
	}
}

func checkStatus(done <-chan any, urls ...string) <-chan Result {
	responses := make(chan Result)
	go func() {
		defer close(responses)
		for _, url := range urls {
			resp, err := http.Get(url)
			select {
			case responses <- Result{resp, err}:
			case <-done:
				return
			}
		}
	}()
	return responses
}
