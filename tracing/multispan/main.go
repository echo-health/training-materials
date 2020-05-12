package main

import (
	"context"
	"fmt"
	"time"

	opentracing "github.com/opentracing/opentracing-go"

	"github.com/echo-health/training-materials/tracer"
)

func main() {
	// Create a tracer - note we don't save a reference to the returned
	// tracer. This is ok as the tracer is stored as the global tracer
	// in opentracing.
	tracer.Tracer("multi-span")

	// Create a new context. The Background() context is empty and
	// never cancelled. When handling an HTTP or gRPC request you
	// can use the context from that request but since we're not
	// inside a request we have to make one.
	ctx := context.Background()

	// Start a new span from the context. This context doesn't contain
	// a reference to an existing span so a "root" span will be returned
	// along with a new context that contains a reference to this span
	span, ctx := opentracing.StartSpanFromContext(ctx, "complex math")

	for _, v := range []int{1, 2, 3} {
		// For each iterartion of the loop create a new span using the
		// context we were given by the earlier call to StartSpanFromContext
		// As this ctx does contain a reference to a span this new span
		// will be a "child" of that one.
		span, _ := opentracing.StartSpanFromContext(ctx, "loop")

		// Set a tag on this span and finish it
		span.SetTag("computed_value", v*2)
		span.Finish()

		time.Sleep(time.Duration(v) * time.Second)
	}

	// Finish the parent span
	span.Finish()

	// Output the trace ID so we can open it in Zipkin
	tracer.LogTraceID(span)

	// Make sure all data has been sent to Zipkin
	time.Sleep(3 * time.Second)
	fmt.Println("Done")
}
