package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/echo-health/zipkin-tutorial/tracer"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/sanity-io/litter"
)

func main() {
	// Create a tracer
	tracer.Tracer("cross-service")

	// Start a new span
	span, _ := opentracing.StartSpanFromContext(context.Background(), "cross service trace")

	// Make a new HTTP request object
	req, err := http.NewRequest("GET", "http://localhost:8001", nil)
	if err != nil {
		panic(err)
	}

	// A carrier is a way of passing opentracing data between different services.
	// In this case we are using a carrier that passes the opentracing data
	// in HTTP headers.
	carrier := opentracing.HTTPHeadersCarrier(req.Header)

	// Add the span data (span id, parent id) to the HTTP headers
	//
	// span.Context() returns a SpanContext - not to be confused with context.Context
	// A SpanContext contains information about the span e.g. it's ID
	//
	// opentracing.HTTPHeaders tells the Inject function what format we want to send the data
	// in. In this case HTTP header format e.g. "key: value"
	err = opentracing.GlobalTracer().Inject(span.Context(), opentracing.HTTPHeaders, carrier)
	if err != nil {
		panic(err)
	}

	// Dump the request headers so we can see what's been injected
	fmt.Println("Request headers:")
	litter.Dump(req.Header)

	// Make the request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	// Tag the span with some information about the response
	span.SetTag("status_code", res.StatusCode)
	span.SetTag("status_message", res.Status)

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	// This adds a log to the span with a key of "response" and a value
	// of the response text.
	//
	// Logs are similar to tags in that they are
	// key->value pairs but unlike tags they cannot be used to search for
	// a span in the Zipkin UI.
	//
	// A span can only have one value for a given tag - newer values overwrite
	// older ones, but you can keep logging the same key multiple times and
	// all will be recorded against the span.
	//
	// More info here https://opentracing.io/docs/overview/tags-logs-baggage/
	span.LogKV("response", string(responseBody))

	// Finish the span
	span.Finish()

	// Output the trace ID so we can open it in Zipkin
	tracer.LogTraceID(span)

	time.Sleep(3 * time.Second)
	fmt.Println("Done")
}
