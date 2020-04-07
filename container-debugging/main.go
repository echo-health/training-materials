package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
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
		time.Sleep(10 * time.Second)
	}
}

func exposePprof() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
}
