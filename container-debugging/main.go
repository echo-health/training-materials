package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
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

func exposePrometheus() {
	go func() {
		// create a new mux server
		server := http.NewServeMux()
		// register a new handler for the /metrics endpoint
		server.Handle("/metrics", promhttp.Handler())
		// start an http server using the mux server
		http.ListenAndServe(":9001", server)
	}()
}
