
# Exercise 1: deployments
1. Run `kubectl config get-contexts` - it's good to start off work by checking you're in the right cluster. you should see output like:
    ```
    charlottegodley@Charlottes-MBP kube-tutorial % kubectl config get-contexts
    CURRENT   NAME                 CLUSTER          AUTHINFO         NAMESPACE
    *         docker-desktop       docker-desktop   docker-desktop   
    ```
1. If it is not set to this, run `kubectl config use-context docker-desktop`. If this fails, you probably don't have the docker desktop integration set up - see the readme.
1. Run `kubectl get ns` - this lists the namespaces you have with vanilla minikube. Not all of these are accessible/modifiable in GKE but would be in vanilla kubernetes.
1. Run `kubectl apply -f manifests/deployment.yaml` from within the location you checked out this repo - you should see:
```
charlottegodley@Charlottes-MBP kube-tutorial % kubectl apply -f manifests/deployment.yaml
namespace/hello-world created
deployment.apps/hello-world created
```
1. Let's open up the file and talk about what's in there. We can inspect the cluster to see what happened here also:
    1. Run `kubectl get ns` and see there's a new namespace
    1. Run `kubectl config set contexts.docker-desktop.namespace hello-world` to set our default query namespace (for the minikube context) to hello-world - optional, but this allows you to skip out `-n hello-world` when poking around with pods.
    1. Run `kubectl get deployment -o yaml` - we see basically the same output as what we put in, except there's new stuff like "status" which tells us what's going on with the deployment, and anything that's system level or which we didn't bother filling in, kubernetes has dumped the default values. Don't worry too much about the stuff in here, a lot of it you'll never need to know.
    1. Running `kubectl describe deployment hello-world` gives us basically the same thing, but with a nicely formatted view.
1. From `describe` you should have noticed that all 3 "replicas" (posh word for "copies") are failing - why is this? (no cheating and looking further down this file)
1. `kubectl get pod` shows us a slightly different view on the world:
    ```
    charlottegodley@Charlottes-MBP kube-tutorial % kubectl get pod
    NAME                           READY   STATUS             RESTARTS   AGE
    hello-world-74b6c95894-j99tk   0/1     CrashLoopBackOff   5          5m47s
    hello-world-74b6c95894-n5872   0/1     CrashLoopBackOff   5          5m47s
    hello-world-74b6c95894-qdt4h   0/1     CrashLoopBackOff   5          5m47s
    ```
    Between our deployment and the pod, kubernetes has assigned a randomised string to the end to make sure we don't have naming clashes.


1. Again let's look at `kubectl describe pod` - if we do this with no pod name, it gives us a full log of allllllll the things. We can talk about some of the fields here if they're interesting, some that in the real world are fun:
    - `Node:         minikube/<vm ip address>`: if you weren't using minikube, you'd hope your pods landed on different nodes. As it stands, we haven't told it to send them anywhere in particular which you can see from `Node-Selectors:  <none>` so Kubernetes magically finds out how much space there is on each node and assigns them according to your requests and limits on memory and storage space.
1. Let's get down to the meat of why this container is sad, at the bottom of the output:
    ```
    Normal   Scheduled  9m23s                   default-scheduler  Successfully assigned hello-world/hello-world-74b6c95894-n5872 to minikube
    Normal   Created    8m30s (x4 over 9m17s)   kubelet, minikube  Created container hello-world
    Normal   Started    8m30s (x4 over 9m17s)   kubelet, minikube  Started container hello-world
    Normal   Pulling    7m49s (x5 over 9m22s)   kubelet, minikube  Pulling image "hello-world:latest"
    Normal   Pulled     7m45s (x5 over 9m17s)   kubelet, minikube  Successfully pulled image "hello-world:latest"
    Warning  BackOff    4m17s (x24 over 9m12s)  kubelet, minikube  Back-off restarting failed container
    ```
    So, we got the container onto a node. The container was created and started...then it failed. Back-off basically means "let's wait a little bit and try it again" (which is fairly common elsewhere in microservices too I suppose)

1. Pick out a pod and run `kubectl logs <podname>` and...
    ```
    Hello from Docker!
    This message shows that your installation appears to be working correctly.

    To generate this message, Docker took the following steps:
    1. The Docker client contacted the Docker daemon.
    2. The Docker daemon pulled the "hello-world" image from the Docker Hub.
        (amd64)
    3. The Docker daemon created a new container from that image which runs the
        executable that produces the output you are currently reading.
    4. The Docker daemon streamed that output to the Docker client, which sent it
        to your terminal.

    To try something more ambitious, you can run an Ubuntu container with:
    $ docker run -it ubuntu bash

    Share images, automate workflows, and more with a free Docker ID:
    https://hub.docker.com/

    For more examples and ideas, visit:
    https://docs.docker.com/get-started/
    ```
    So, it turns out we're running the wrong container, or don't have it locally. Good to note: if this was a pod controlled by a CronJob, this wouldn't be a "fail", just a "completed". Pods ran by deployments are expected to never end execution.

1. So...let's actually build the container and tag it as it is in the deployment file. `docker build -t hello-world:latest .`
1. `kubectl get pod`...and...nothing...happened?
    ```
    charlottegodley@Charlottes-MBP kube-tutorial % kubectl get pod
    NAME                           READY   STATUS             RESTARTS   AGE
    hello-world-74b6c95894-j99tk   0/1     CrashLoopBackOff   12         37m
    hello-world-74b6c95894-n5872   0/1     CrashLoopBackOff   12         37m
    hello-world-74b6c95894-qdt4h   0/1     CrashLoopBackOff   12         37m
    ```
1. Let's go back to the deployment.yaml file and look at what image we set - `image: hello-world:latest`. Hmm, matches the tag we built.
...
1. The problem here is `imagePullPolicy` - when we first started the pod, it was pulled from the docker registry. This means we already _have_ a copy of hello-world:latest and it's not our copy. The default policy is `ifNotPresent` meaning, if there's no image, pull it down. 
    - Something else to note here: if we switch this to `Always` (which would, after running `kubectl delete pod --all`), it's still bad practice to use `latest` tag because your code (whatever version it is since you can't tell from the tag) could randomly get deployed.
1. In this context (i.e minikube) we need to set it to `Never`, because we've mapped the docker daemon on the machine to the one we used to build the image, so...we should _always_ have the image there anyway, so don't try to pull it from docker.
    ```
            spec:
                containers:
                - name: hello-world
                    image: hello-world:1.0
                    imagePullPolicy: Never
                    ports:
                    - containerPort: 8080 # need to specify this so that it is exposed to other things in the namespace e.g services
    ```
    - we've also changed the tag to be `1.0.0` to get around using latest so will need to retag/rebuild: `docker build -t hello-world:1.0.0 .`
1. run `kubectl apply -f manifests/deployment.yaml` and it should work.
    ```
    charlottegodley@Charlottes-MBP kube-tutorial % kubectl get pod
    NAME                           READY   STATUS    RESTARTS   AGE
    hello-world-675bc74777-f4rkr   1/1     Running   0          3m34s
    hello-world-675bc74777-jb5n5   1/1     Running   0          3m32s
    hello-world-675bc74777-tmq9v   1/1     Running   0          3m33s
    ```
    1. In another window, you can run `kubectl get pod -w` to watch what kubernetes does in this situation - you should see the old pods `Terminating` one by one as new pods go into `ContainerCreating` and eventually `Running`
1. Now, up until now, we've just looked at the state of the pod, deployment etc without looking at what it's doing...but it's inside a cluster, inside a VM, so from where we're running these commands, we can't actually send any traffic to it. To do this we use `port-forward` which works similarly to mounting a docker port to your computer's ports: `kubectl port-forward <pod-name> 8080:8080`
1. Now if you open up a web browser and go to `http://localhost:8080/world`, you should see `Hello, world!`.