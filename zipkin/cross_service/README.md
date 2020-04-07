# Cross-service Example

This is an example of how a trace can cover multiple services, much like how we do in production.

### How to run

- Start Zipkin. Either `docker-compose up` or see [running in Kubernetes](../kubernetes)
- Open the [Zipkin UI](http://localhost:9411/zipkin/)
- Run the "other" service - `go run other_service/main.go`
- Run (in a different tab) the Go program - `go run main.go`
- Look at the output from each program to see how trace information was passed between services
- Find the span in Zipkin

### What is happening?

In this example service A (main.go) makes an HTTP request to service B (other_service/main.go). To make sure that the spans created in service B can be linked to the one created in service A, service A sends information about the current span in the HTTP headers. Service B then extracts that information from the received headers and creates a span using that information.

In this example the two services are communicating over HTTP but the communication method doesn't really matter. For example in production we have services that communicate over gRPC, HTTP, and pubsub. The thing to understand is that by passing information about the current span to another service or process it can continue to add spans to underlying trace. This is what is meant by cross-service tracing.

Things you can try after running this example:

- Try making multiple requests to service B from service A and see how that shows up in Zipkin
- Try sending some data to service B in query parameters and then logging that data to a span in service B
