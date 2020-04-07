package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
)

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

func exposePprof() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
}
