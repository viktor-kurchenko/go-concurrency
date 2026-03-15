package main

import "fmt"

func main() {
	fmt.Printf("Error handling pattern example\n")
	errorHandlingExample()

	fmt.Printf("\nPipeline pattern example\n")
	pipelineExample()

	fmt.Printf("\nFanOut FanIn pattern example\n")
	fanOutFanInExample()

	fmt.Printf("\nOrDone pattern example\n")
	orDoneExample()

	fmt.Printf("\nTee pattern example\n")
	teeExample()
}
