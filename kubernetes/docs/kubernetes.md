# Precursor: What is Kubernetes?
Most devs are fairly familiar with **containers**. What do we use those for? We use them to encapsulate everything that goes into an application - what language is needed, what dependent packages need to be installed, so that we can throw the app at a server and run it somewhere.

The *somewhere* would be a random docker run command on a server, but the problem with that is:
1. You would need to set that server up yourself and monitor it
1. You would need to monitor docker and your application to make sure it's up all the time
1. If you want it to be "highly available" (i.e, have multiple copies of your application accross multiple data centers) you would need to manage the additional servers and do some load balancing to make the traffic distribute across those copies.

A **container orchestrator** handles many of these problems. Kubernetes is one of those, as well as **mesos, docker swarm** and probably more. The problem with most orchestrators is they might literally just know how to run and restart containers - Kubernetes also covers networking (including some load balancing), security, access control and gives you the ability to further extend it.

## Where does the boundary lie between a cloud provider and kubernetes?
Kubernetes is generally concerned with the node downwards, i.e it does not handle creating machines, deleting machines, auto scaling up to more nodes.

Cloud providers handle the node upwards. You can install Kubernetes on any cloud provider or set of machines whether or not they have their own integrations. When creating a GKE environment for example, you will configure what type of machine to use and how many of them.

Where the boundary is blurred slightly:
- Ingress: the concept of traffic coming into the cluster. Many cloud providers will provide an abstraction within their managed kubernetes offering to create things like static IP addresses on the fly.
- Storage: if you are rolling your own cluster, you will likely need to install extensions to plug it into persistent storage for the cloud provider or environment you are using.
