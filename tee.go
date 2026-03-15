package main

import (
	"fmt"
)

func teeExample() {
	input := make(chan int)
	go func() {
		defer close(input)
		for i := range 5 {
			input <- i
		}
	}()

	done := make(chan any)
	defer close(done)

	output1, output2 := tee(done, input)
	for out1 := range output1 {
		fmt.Printf("output1: %d, output2: %d\n", out1, <-output2)
	}
}

func tee(done <-chan any, input <-chan int) (<-chan int, <-chan int) {
	output1 := make(chan int)
	output2 := make(chan int)
	go func() {
		defer close(output1)
		defer close(output2)
		for i := range orDone(done, input) {
			out1, out2 := output1, output2
			for range 2 {
				select {
				case out1 <- i:
					out1 = nil
				case out2 <- i:
					out2 = nil
				case <-done:
					return
				}
			}
		}
	}()
	return output1, output2
}
