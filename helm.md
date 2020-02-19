# Exercise 4: Helm
In this exercise we're going to learn how to create and deploy a helm chart.

1. [Install helm](https://helm.sh/docs/intro/install/) (TODO: make this a pre-req)
1. Run `helm create hello-world` in a terminal in this directory. This creates what's called a "chart" which means "templated kubernetes yaml". 
1. Build the app to a semver tag - `docker build -t hello-world:1.0.0 .`
1. Open up `hello-world/templates/deployment.yaml` and `hello-world/values.yaml` and configure values so that it will deploy hello-world version 1.0.0 (figure this out and discuss it)
1. Now we're going to deploy it - `helm install hello-world hello-world` (the first argument is the name within the cluster, the second is the chart's name)
1. Now we can do what we did before to debug it - check the deployment, check the pods, check the app is working using port-forward.
1. Make a change to the application (whatever you like, e.g, change the output) and rebuild it - `docker build -t hello-world:1.0.1 .`, then redeploy it. You can do this by upping the app version (discover how to do this for yourself) then running `helm upgrade hello-world hello-world` 
1. Check that's worked - `kubectl get pod`
