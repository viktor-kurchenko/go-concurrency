package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

func fanOutFanInExample() {
	done := make(chan any)
	defer close(done)
	for i := range fanIn(done, fanOut(done, 10)...) {
		fmt.Printf("result: %d\n", i)
	}
}

func fanOut(done <-chan any, count int) []<-chan int {
	outputs := make([]<-chan int, 0, count)
	for range count {
		outputs = append(outputs, longTimeCalculation(done))
	}
	return outputs
}

func fanIn(done <-chan any, inputs ...<-chan int) <-chan int {
	output := make(chan int)
	go func() {
		wg := sync.WaitGroup{}
		for _, i := range inputs {
			wg.Add(1)
			go func(in <-chan int) {
				defer wg.Done()
				for i := range in {
					select {
					case output <- i:
					case <-done:
						return
					}
				}
			}(i)
		}

		go func() {
			defer close(output)
			wg.Wait()
		}()
	}()
	return output
}

func longTimeCalculation(done <-chan any) <-chan int {
	output := make(chan int)
	go func() {
		defer close(output)
		select {
		case <-time.After(time.Second):
			output <- rand.Int()
		case <-done:
		}
	}()
	return output
}
