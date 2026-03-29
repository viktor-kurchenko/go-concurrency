package main

import (
	"fmt"
	"time"
)

func keepalive() {
	done := make(chan any)
	defer close(done)

	timeout := time.Second * 2
	result, keepalive := startJob(done, timeout)

	for range 15 {
		select {
		case r, ok := <-result:
			if !ok {
				return
			}
			fmt.Printf("result: %d\n", r)
		case k, ok := <-keepalive:
			if !ok {
				return
			}
			fmt.Printf("keepalive: %v\n", k)
		case <-time.After(timeout):
			return
		}
	}
}

func startJob(done <-chan any, timeout time.Duration) (<-chan int, <-chan time.Time) {
	results := make(chan int)
	keepAlive := make(chan time.Time)
	go func() {
		defer close(results)
		defer close(keepAlive)

		pulse := time.Tick(timeout / 4)
		workGen := time.Tick(timeout / 2)

		sendPulse := func() {
			select {
			case keepAlive <- time.Now():
			default:
			}
		}

		workIdx := 0
		for {
			select {
			case <-workGen:
				results <- workIdx
				workIdx++
			case <-pulse:
				sendPulse()
			case <-done:
				return
			}
		}
	}()
	return results, keepAlive
}
