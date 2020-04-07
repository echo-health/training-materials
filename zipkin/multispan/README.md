# Multi-span Example

This is an example of how to create child spans.

### How to run

- Start Zipkin. Either `docker-compose up` or see [running in Kubernetes](../kubernetes)
- Open the [Zipkin UI](http://localhost:9411/zipkin/)
- Run the Go program - `go run main.go`
- Find the span in Zipkin

### What is happening?

This example creates a root (or parent) span that has a number of child spans. Each child span can have it's own name, tags and logs. This example shows how you can use spans to capture information at a granular level in applications. Things you can try after running this example:

- Try adding more iterations to the loop
- Try adding another loop inside the existing one that creates yet another child span. See how deep you can go...
