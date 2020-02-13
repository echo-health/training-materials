# kubernetes-tutorial
This repository is intended to teach what the abstractions are in Kubernetes, how to use them and how to debug them. Each exercise is self contained, but it is recommended to start at exercise 1 if you are a beginner.

## Prerequisites
- You need to set up your mac to use the docker-for-mac kubernetes integration: [instructions](https://docs.docker.com/docker-for-mac/#kubernetes). This will spin up a kubernetes cluster inside a container and comes with kubectl in case you don't already have this.
- For one bit you will also need [jq](https://stedolan.github.io/jq/) which you can get by running `brew install jq`

## Exercises
- [Exercise 1: Deployments](deployments.md)
- [Exercise 2: Services](services.md)
- [Exercise 3: Config maps and secrets](config.md)

## Wider reading, references etc
- The kubernetes documentation is fairly heavy. [kubernetes.io](https://kubernetes.io) - Generally the best way to use this is to google the resource type you're interested in
- https://github.com/ramitsurana/awesome-kubernetes - list of other "awesome" resources for learning more about kubernetes.