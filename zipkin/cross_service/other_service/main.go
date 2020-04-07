package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/echo-health/zipkin-tutorial/tracer"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"

	"github.com/sanity-io/litter"
)

func main() {
	tracer.Tracer("other-service")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Dump the HTTP headers so we can see what was sent
		fmt.Println("HTTP Headers:")
		litter.Dump(r.Header)

		// Exactly the same as when we made an HTTP request, when receiving
		// an HTTP request we can make a carrier from the headers
		carrier := opentracing.HTTPHeadersCarrier(r.Header)

		// Extract the opentracing data and use it to make a SpanContext
		parentSpanContext, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, carrier)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		// Now we can start a new span that is a child of the parent span (from the calling service).
		span := opentracing.GlobalTracer().StartSpan("doing stuff", opentracing.ChildOf(parentSpanContext))

		// So far we've been calling span.Finish() when we're done with the span.
		// This is fine but in Go programs it's common to return from a function
		// many times e.g. if err != nil { return err }. As it would be boring
		// and error-prone to have to remember to call span.Finish() in all these
		// places we can use a defer to make sure the span is always finished after
		// this function returns.
		// More on defer here - https://tour.golang.org/flowcontrol/12
		defer span.Finish()

		// Log some information to the span. There are few different functions
		// for logging to a span. This one is more type-safe than using LogKV
		// as you have to explicitly state the types of each field.
		span.LogFields(log.String("name", "Jon"), log.Int("age", 33))

		// So that we can see a span in Zipkin that has a meaningful duration we'll
		// sleep for a bit here. In a normal app you'd probably do things like read
		// or write from a database or make requests to other services, all things that
		// that take time.
		time.Sleep(1 * time.Second)

		// Send a response back to the client
		w.Write([]byte("Hello from other service"))
	})

	fmt.Println("Running...")
	err := http.ListenAndServe("0.0.0.0:8001", http.DefaultServeMux)
	if err != nil {
		panic(err)
	}

	fmt.Println("Done")
}
