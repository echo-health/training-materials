# Simple Example

This is possibly the simplest example of sending some tracing data to Zipkin.

### How to run

- Start Zipkin. Either `docker-compose up` or see [running in Kubernetes](../kubernetes)
- Open the [Zipkin UI](http://localhost:9411/zipkin/)
- Run the Go program - `go run main.go`
- Find the span in Zipkin

### What is happening?

This example sends a single span to Zipkin containing a single tag. Things you can try after running this example:

- Try changing the tracer name and re-running
- Try adding different tags and see how they show up in the UI
