package tracer

import (
	"flag"
	"fmt"
	"log"
	"os"

	opentracing "github.com/opentracing/opentracing-go"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
)

var port string

func init() {
	flag.StringVar(&port, "tracingPort", "9411", "port the tracer is running on")
	flag.Parse()
}

func Tracer(name string) opentracing.Tracer {

	// This method uses the Zipkin client libraries to create a tracer even though we are using
	// Jaeger as a backend. This is ok as Zipkin and Jaeger are compatible with each other.

	// An endpoint represents a specific service
	// This can be a host/port combo or it can just be the name
	// of the service (which is what we do)
	endpoint, err := zipkin.NewEndpoint(name, "")
	if err != nil {
		panic(err)
	}

	// This is where the tracing data will be sent. We're running Jaeger
	// locally so it's available on localhost but in Kubernetes it's
	// available using Kubernetes' DNS e.g. http://tracing/
	addr := fmt.Sprintf("http://localhost:%s/api/v2/spans", port)

	// Create a reporter to send data to Jaeger. We use an HTTP reporter
	// which sends data to Jaeger over HTTP requests.
	// There are other ways of sending data to zipkin e.g. RabbitMQ, gRPC
	reporter := zipkinhttp.NewReporter(addr, zipkinhttp.Logger(log.New(os.Stdout, "TRACING: ", log.LstdFlags)))

	// Create a tracer using our reporter
	nativeTracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint))
	if err != nil {
		log.Fatalf("unable to create tracer: %+v\n", err)
	}

	// This is so we can use the Zipkin tracer as an "opentracing" tracer
	tracer := zipkinot.Wrap(nativeTracer)

	// Set this tracer as the global tracer. This means it will be returned
	// by opentracing.GlobalTracer() or used when you call opentracing.StartSpan()
	opentracing.SetGlobalTracer(tracer)

	return tracer
}

func LogTraceID(span opentracing.Span) {
	zs, ok := span.Context().(zipkinot.SpanContext)
	if ok {
		traceID := zs.TraceID.String()
		fmt.Printf("View trace: http://localhost:%s/trace/%s\n", port, traceID)
	}
}
