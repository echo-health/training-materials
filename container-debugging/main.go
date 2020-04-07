package main

import "fmt"

func main() {
	channel := make(chan int)
	count := 0
	for {
		count += 1
		fmt.Printf("Starting a new goroutine - count %d\n", count)
		go func() {
			<-channel
		}()
	}
}
