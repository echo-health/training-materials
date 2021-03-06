# Running Jaeger in Kubernetes

In this directory is the minimum config required for running Jaeger in a local Kubernetes cluster.

## How to

- Install [Docker for Mac](https://download.docker.com/mac/stable/Docker.dmg) and [set up Kubernetes](https://docs.docker.com/docker-for-mac/#kubernetes)
- Apply the manifest in this directory - `kubectl apply -f ./manifest.yaml`
- Run `kubectl get service -n tracing-tutorial` to see which port Jaeger is running on (see below)
- Add `-tracingPort <PORT> -webPort <PORT>` to all `go run` commands in the examples
- When you're done run `kubectl delete ns/tracing-tutorial` to stop running Jaeger in Kubernetes

### Finding port Jaeger is running on

The output of `kubectl get service -n tracing-tutorial` should be something like:

```
$ kubectl get service -n tracing-tutorial

NAME     TYPE       CLUSTER-IP      EXTERNAL-IP   PORT(S)                          AGE
jaeger   NodePort   10.106.66.192   <none>        9411:30486/TCP,16686:30452/TCP   8s
```

In this example you would need to add `-tracingPort=30486 -webPort=30452` to all the `go run` commands in the examples. This because the tracing API is exposed on `9411` (which is mappped tp `30486`) and the web UI is exposed on `16686` (which is mapped to `30452`) Everytime you destroy and re-create the resources in Kubernetes you'll get a new port.
