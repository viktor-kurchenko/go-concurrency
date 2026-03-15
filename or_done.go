package main

import (
	"fmt"
)

func orDoneExample() {
	input := make(chan int)
	go func() {
		defer close(input)
		for i := range 5 {
			input <- i
		}
	}()
	done := make(chan any)
	defer close(done)
	for i := range orDone(done, input) {
		fmt.Printf("value: %d\n", i)
	}
}

func orDone(done <-chan any, input <-chan int) <-chan int {
	output := make(chan int)
	go func() {
		defer close(output)
		for {
			select {
			case v, ok := <-input:
				if !ok {
					return
				}
				select {
				case <-done:
				case output <- v:
				}
			case <-done:
			}
		}
	}()
	return output
}
