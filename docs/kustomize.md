# Exercise 5: Kustomize
We're now going to try doing the same thing as before with kustomize.

1. Take a look at the files in `kustomize` - note that the bases can be used on their own, or with overlays to shape them slightly for one environment or another
1. Apply it to your cluster using `kubectl apply -k kustomize/overlays/production`


The intention is to move to kustomize over helm at echo to reduce some of the complexity.