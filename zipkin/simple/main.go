package main

import (
	"fmt"
	"time"

	"github.com/echo-health/zipkin-tutorial/tracer"
)

func main() {
	// Create a tracer
	t := tracer.Tracer("simple")

	// Create a new span. This span has no parent so
	// is referrred to as a "root" span.
	span := t.StartSpan("simple span")

	// Add a tag to a span. Tags can make it easy to
	// find a specific trace e.g. you could add a tag like:
	//     span.SetTag("patient_id", patient.ID)
	// Then in Zipkin you can search by that tag
	span.SetTag("foo", "bar")

	// Finish the span. This will cause the data for
	// this span to be sent to Zipkin (using the tracer
	// we made earlier)
	span.Finish()

	// Output the trace ID so we can open it in Zipkin
	tracer.LogTraceID(span)

	// Spans are not sent immediately but in batches every few seconds.
	// In a real world app we wouldn't need this as the app keeps
	// running but to make sure the span gets sent in this
	// example we add a small sleep.
	time.Sleep(3 * time.Second)
	fmt.Println("Done")
}
