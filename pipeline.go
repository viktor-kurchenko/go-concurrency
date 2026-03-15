package main

import "fmt"

func pipelineExample() {
	done := make(chan any)
	defer close(done)
	for i := range multiply(done, increment(done, intStream(done, 10), 1), 2) {
		fmt.Printf("result: %d\n", i)
	}
}

func intStream(done <-chan any, count int) <-chan int {
	stream := make(chan int)
	go func() {
		defer close(stream)
		for i := 1; i <= count; i++ {
			select {
			case stream <- i:
			case <-done:
				return
			}
		}
	}()
	return stream
}

func increment(done <-chan any, input <-chan int, inc int) <-chan int {
	output := make(chan int)
	go func() {
		defer close(output)
		for i := range input {
			select {
			case output <- i + inc:
			case <-done:
				return
			}
		}
	}()
	return output
}

func multiply(done <-chan any, input <-chan int, mul int) <-chan int {
	output := make(chan int)
	go func() {
		defer close(output)
		for i := range input {
			select {
			case output <- i * mul:
			case <-done:
				return
			}
		}
	}()
	return output
}
