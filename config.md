# Exercise 3: Configuring Applications
1. As before, check we are in the right cluster: `kubectl config get-contexts` and where appropriate, set your query namespace to be where we'll be working: `kubectl config set contexts.minikube.namespace hello-world`
1. Open up `manifests/configmap.yaml` - let's discuss what's in here:
    - We have 2 applications in 2 namespaces
    - We also have 2 ways of writing a config map - one uses pure yaml, the other uses yaml strings
    - One application uses the config map to set it's environment variables as shown in yaml, the other mounts it as a file
1. Before we actually use the config maps, we need to change a few things:
    - run `git checkout env`
    - rebuild and tag it as `env` - `docker build -t hello-world:env .`
    - run `git checkout env-file`
    - rebuild and tag it as `env-file` - `docker-build -t hello-world:env-file .`
1. Finally, let's apply the yaml: `kubectl apply -f manifests/configmap.yaml`
1. We can now check what's going on in prod: 
    - `kubectl get pod -n hello-world` should show it's running
    - then `kubectl port-forward <pod-name> 8080:8080` should show `Hello, production!` when we pull this up in a web browser. 
    - Open another tab and either edit the file and apply it again to change the value in the config map to something else, or run `kubectl edit configmap app-config` to edit it directly on the cluster. Now refresh the web browser and...err...nothing happened :D
    - This is because loading config maps in kubernetes by environment variable doesn't refresh the config in the pod, or cause it to die and restart. So, to get the change to be applied correctly, you're going to need to delete the pods :D
1. And staging:
    - `kubectl get pod -n hello-world-staging` should show it _is not_ running :D
    - to debug this, let's use exec: `kubectl exec <pod name> -it -- /bin/sh` (`-it` means interactive, `-- /bin/sh` is what command we are running)
    - and run `ls /etc/config` to find out the filename we mounted the config map as.
    - It appears to be named "environment" because that's what we wrote in the config map. Open up `main.go` and figure out what's wrong, and suggest how we fix this :D