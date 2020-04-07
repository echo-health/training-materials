# Exercise 4: Helm
In this exercise we're going to learn how to create and deploy a helm chart.

1. Run `helm create hello-world` in a terminal in this directory. This creates what's called a "chart" which means "templated kubernetes yaml". 
1. Build the app to a semver tag - `docker build -t hello-world:1.0.0 .`
1. Open up `hello-world/templates/deployment.yaml` and `hello-world/values.yaml` and configure values so that it will deploy hello-world version 1.0.0 (figure this out and discuss it)
1. Now we're going to deploy it - `helm install hello-world hello-world` (the first argument is the name within the cluster, the second is the chart's name)
1. Now we can do what we did before to debug it - check the deployment, check the pods, check the app is working using port-forward.
1. Make a change to the application (whatever you like, e.g, change the output) and rebuild it - `docker build -t hello-world:1.0.1 .`, then redeploy it. You can do this by upping the app version (discover how to do this for yourself) then running `helm upgrade hello-world hello-world` 
1. Check that's worked - `kubectl get pod`

This is, in a loose fashion, how our CI system currently works - we have a similar looking chart in [this repo](https://github.com/echo-health/infrastructure/tree/master/platform/charts/service) which is the same for every service. The `values.yaml` file is then overriden on a per-service basis to configure the chart closer to where the code lives.

Jenkins takes the charts and does some magic to deploy them using [this groovy script](https://github.com/echo-health/jenkins-library/blob/master/vars/helm.groovy).
