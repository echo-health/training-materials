# Tracing Training

This directory contains some examples of how to use tracing to instrument applications.

### What is instrumentation?

At a high level instrumentation is adding things to your code that allow you to debug it's performance or diagnose errors.

Logging could be considered a type of instrumentation although in a large distributed system (one containing many parts/services) logging quickly loses it's value. This is because trying to stitch together logs from multiple services where each service is handling many requests at a time can be quite tricky.

[Jaeger](https://www.jaegertracing.io/) is a cross-service (or distributed) tracing system. It was originally developed by Uber to help them debug their services architecture.

### Setup

- Clone this repo
- [Docker for Mac](https://download.docker.com/mac/stable/Docker.dmg)
- Optional - Docker For Mac kubernetes setup - [instructions](https://docs.docker.com/docker-for-mac/#kubernetes). This will spin up a kubernetes cluster inside a container and comes with kubectl in case you don't already have this. (If you do already have `kubectl`, it should Just Work)

#### Check Go environment

To check your Go envionment run `go get ./...` and `go vet ./...` inside the repo after cloning. If either of these commands error you may not have a working Go environment. If so seek help in Slack in the #engineering or #platform channels.

#### Check Docker environment

After cloning run `docker-compose up` inside the repo. This should first download the [Jaeger Docker image](https://hub.docker.com/r/jaegertracing/all-in-one) and then start Jaeger. You should be able to open the [Jaeger UI](http://localhost:16686/) in a web browser.

If you get any errors when trying to run this command then you may not have a working Docker environment. If so seek help in Slack in the #engineering or #platform channels.

You can use `Ctrl + C` to stop the Jaeger process from running. This will just stop the process, if you want to remove it completely then use `docker-compose rm`.

#### Check Kubernetes environment (optional)

If you want to run Jaeger locally using Kubernetes it's worth checking that your environment is in the right state. To check this run:

```bash
$ kubectl config get-contexts
```

You should see a list of available "contexts", one of which should be `docker-desktop`. If you don't see this or get any errors ask for help in the #engineering or #platform Slack channels.

If the `docker-desktop` context doesn't have a `*` next to it under the `CURRENT` column then run `docker config use-context docker-desktop` to use that context for this tutorial. The follow the [instructions for running Jaeger in Kubernetes](kubernetes).

### Workshop

This repo can be used to run a short workshop on using tracing. A suggested agenda for this type of workshop would be:

- What is tracing?
- The difference between [OpenTracing](https://opentracing.io/docs/overview/what-is-tracing/) and [Jaeger](https://www.jaegertracing.io/)
- Example 1: [simple](simple)
- Example 2: [multispan](multispan)
- Example 3: [cross_service](cross_service)
- Examine a [real-world production trace](https://trace.echo.co.uk) - the `orders/CreateOrders` RPC has quite a lot going on 🤓
- Questions
