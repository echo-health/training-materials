# kubernetes-tutorial
This repository is intended to teach what the abstractions are in Kubernetes, how to use them and how to debug them. Each exercise is self contained, but it is recommended to start at exercise 1 if you are a beginner.

## Prerequisites
- Workshop 1:
    - [Docker for Mac](https://download.docker.com/mac/stable/Docker.dmg)
    - You need to set up your mac to use the Docker For Mac kubernetes integration: [instructions](https://docs.docker.com/docker-for-mac/#kubernetes). This will spin up a kubernetes cluster inside a container and comes with kubectl in case you don't already have this. (If you do already have `kubectl`, it should Just Work)
    - You will also need [jq](https://stedolan.github.io/jq/) which you can get by running `brew install jq`
- Workshop 2: 
    - [helm](https://helm.sh/docs/intro/install/)

## Exercises
### Workshop 1: Vanilla Kubernetes
- [(Reading) Precursor: what is kubernetes?](kubernetes.md)
- [Exercise 1: Deployments](deployments.md)h
- [Exercise 2: Services](services.md)
- [Exercise 3: Config maps and secrets](config.md)

### Workshop 2: Echo-Flavoured Kubernetes
- [Exercise 4: Helm](helm.md)
- [Exercise 5: Kustomize](kustomize.md)

## Wider reading, references etc
- The kubernetes documentation is fairly heavy. [kubernetes.io](https://kubernetes.io) - Generally the best way to use this is to google the resource type you're interested in
- https://github.com/ramitsurana/awesome-kubernetes - list of other "awesome" resources for learning more about kubernetes.