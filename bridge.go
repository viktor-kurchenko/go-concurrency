package main

import "fmt"

func bridgeExample() {
	input := make(chan (<-chan int))
	go func() {
		defer close(input)
		for i := range 5 {
			ch := make(chan int, 1)
			ch <- i
			close(ch)
			input <- ch
		}
	}()

	done := make(chan any)
	defer close(done)

	for i := range bridge(done, input) {
		fmt.Printf("value: %d\n", i)
	}
}

func bridge(done <-chan any, input <-chan <-chan int) <-chan int {
	output := make(chan int)
	go func() {
		defer close(output)
		for {
			var ch <-chan int
			select {
			case tmpCh, ok := <-input:
				if !ok {
					return
				}
				ch = tmpCh
			case <-done:
				return
			}
			for i := range orDone(done, ch) {
				select {
				case output <- i:
				case <-done:
					return
				}
			}
		}
	}()
	return output
}
