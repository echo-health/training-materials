# Running Zipkin in Kubernetes

In this directory is the minimum config required for running Zipkin in a local Kubernetes cluster.

## How to

- Install [Docker for Mac](https://download.docker.com/mac/stable/Docker.dmg) and [set up Kubernetes](https://docs.docker.com/docker-for-mac/#kubernetes)
- Apply the manifest in this directory - `kubectl apply -f ./zipkin.yaml`
- Run `kubectl get service -n zipkin-tutorial` to see which port Zipkin is running on (see below)
- Add `-zipkinPort <PORT>` to all `go run` commands in the examples
- When you're done run `kubectl delete ns/zipkin-tutorial` to stop running Zipkin in Kubernetes

### Finding port Zipkin is running on

The output of `kubectl get service -n zipkin-tutorial` should be something like:

```
$ kubectl get service -n zipkin-tutorial

NAME     TYPE       CLUSTER-IP      EXTERNAL-IP   PORT(S)          AGE
zipkin   NodePort   10.103.204.66   <none>        9411:30235/TCP   9s
```

In this example Zipkin is running at `http://localhost:30235`. Everytime you destroy and re-create the resources in Kubernetes you'll get a new port.
