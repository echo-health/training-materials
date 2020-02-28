# kubernetes-tutorial
This repository is intended to teach what the abstractions are in Kubernetes, how to use them and how to debug them. Each exercise is self contained, but it is recommended to start at exercise 1 if you are a beginner.

If you are intending on running this as a workshop, I've made a brief [agenda doc](workshop.md) to suggest things to mention during the session.

## Prerequisites
- Workshop 1:
    - Clone this repo, and start a terminal window in the folder where you cloned it
    - [Docker for Mac](https://download.docker.com/mac/stable/Docker.dmg)
    - You need to set up your mac to use the Docker For Mac kubernetes integration: [instructions](https://docs.docker.com/docker-for-mac/#kubernetes). This will spin up a kubernetes cluster inside a container and comes with kubectl in case you don't already have this. (If you do already have `kubectl`, it should Just Work)
    - You will also need [jq](https://stedolan.github.io/jq/) which you can get by running `brew install jq`
- Homework: 
    - [helm](https://helm.sh/docs/intro/install/)

## Exercises
### Workshop 1: Vanilla Kubernetes
- [(Reading) Precursor: what is kubernetes?](docs/kubernetes.md)
- [Exercise 1: Deployments](docs/deployments.md)
- [Exercise 2: Services](docs/services.md)
- [Exercise 3: Config maps and secrets](docs/config.md)

## Workshop 2: Further container debugging
COMING SOON :D

### Homework: How we deploy to Kubernetes at Echo
- [Exercise 4: Helm](docs/helm.md)
- [Exercise 5: Kustomize](docs/kustomize.md)

## Wider reading, references etc
- The kubernetes documentation is fairly heavy. [kubernetes.io](https://kubernetes.io) - Generally the best way to use this is to google the resource type you're interested in
- Aliasing options:
    - [kubectl-aliases](https://github.com/ahmetb/kubectl-aliases) a script which generates a lot of aliases for you so you can do less typing

- [katacoda](https://katacoda.com) is an interactive learning environment for DevOps which has a lot of courses on Kubernetes among other useful tools.
- https://github.com/ramitsurana/awesome-kubernetes - list of other "awesome" resources for learning more about kubernetes.
- Kustomize:
    - [kubernetes.io docs](https://kubectl.docs.kubernetes.io/pages/examples/kustomize.html)
    - [kustomize full documentation](https://kustomize.io)
    - [kustomize repo](https://github.com/kubernetes-sigs/kustomize/tree/master/docs)
- Helm:
    - [helm docs](https://helm.sh)
    - [simple hello world tutorial](https://medium.com/@pablorsk/kubernetes-helm-node-hello-world-c97d20437abd)