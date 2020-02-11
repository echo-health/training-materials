# Exercise 2: Services
1. As before, check we are in the right cluster: `kubectl config get-contexts` and where appropriate, set your query namespace to be where we'll be working: `kubectl config set contexts.minikube.namespace hello-world`
1. Open up `manifest/service.yaml` - this is largely the same as deployment.yaml but with our service added to the top. 
1. Deploy it to the cluster: `kubectl apply -f manifests/service.yaml`
1. Describe it: `kubectl describe service` - what do we see?
    - In particular, we have this line: 
        `Endpoints:         172.17.0.7:8080,172.17.0.8:8080,172.17.0.9:8080`
        this _should_ match the IP addresses of our pods, which we can check with jq:
        `kubectl get pod -o json -l foo=bar | jq -r '.items | .[] | .status | .podIP '`
        and the response (in my case):
        ```
        172.17.0.9
        172.17.0.8
        172.17.0.7
        ```
    - If a service _does not_ have endpoints, for example if you run `kubectl delete deployment hello-service`, this is _normally_ a good indicator that:
        - your application wasn't deployed properly (i.e, there is no deployment and therefore no pods)
        - your labels and selectors aren't correct (i.e, the labels on your pod do not match the labels on the service)
        - something is wrong with your cluster (i.e, one of the system level kubernetes components has a problem and is not updating the link between services and pods). **_This one would be Google's problem in GKE_**
1. To further examplify this, let's change a couple of things:
    1. Give our original application the same label so that it gets picked up by the service:
        ```
          template: # from here on is effectively the pod's yaml
            metadata:
            labels:
                app: hello-world
                foo: bar
        ```
        and then apply it `kubectl apply -f manifests/deployment.yaml`
    1. Update our application so that it doesn't do the exact same thing as the other deployment:
        - Open up `main.go` and change line 14 so that it says `Good morning,` instead of `Hello,`
        - Build it under a new tag: `docker build -t hello-world:good-morning .`
        - Update our deployment in `manifests/service.yaml` to use the new tag:
            ```
                    spec:
                        containers:
                        - name: hello-world
                          image: hello-world:good-morning
                          ports:
                          - containerPort: 8080 
            ```
        - Apply it `kubectl apply -f manifests/service.yaml`
    1. Finally, test it! `kubectl get svc` should now show 6 endpoints. We can also port-forward again: `kubectl port-forward svc/hello-world 8080:80` - note that in testing I never once got the good morning endpoint despite testing it separately so...your mileage may vary.
1. We can also make it more reliable to route traffic to the `hello-service` by changing some more things:
    - update `main.go` so that it runs on port `8081`
    - Change the service so that there's a second port definition pointing to port 8081
    - rebuild, change the tag, update the yaml, reapply. (It should be fairly obvious how to do this part by now :D )
    - **NOTE** you probably wouldn't do it this way in the real world - we would instead add the container as a second container definition on the first deployment. What's the problem in the 2 deployment approach?
1. Finally, we can debug whether the service works _inside_ the cluster by deploying a new pod: `kubectl run curl --image=radial/busyboxplus:curl -i --tty --rm`
    - Then run `nslookup hello-service` to check it's discoverable:
      ```
      Server:    10.96.0.10
      Address 1: 10.96.0.10 kube-dns.kube-system.svc.cluster.local

      Name:      hello-service
      Address 1: 10.103.232.123 hello-service.hello-world.svc.cluster.local
      ```
    - and curl it: `curl hello-service/world`:
      ```
      [ root@curl-crg-5b9787c7bf-8w7fm:/ ]$ curl hello-service/world
        Good morning, world!
      ```
