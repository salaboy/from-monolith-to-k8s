# From Monolith to K8s

Workshop-style guide for moving from a monolith application to a cloud-native architecture running in Kubernetes.

This guide will take you through an example scenario to move from a Monolith application to a set of Cloud Native microservices running in a Kubernetes Cluster. This workshop highlights the use of certain tools to solve some particular challenges that you will face while going to the cloud. 

# Prerequisites

- Kubernetes Cluster
  - Tested in GKE Cluster (4 nodes - n1-standard-2)
- `kubectl` configured. 
- 

# Tools that we are using
- [Kubernetes](http://kubernetes.io)
- [Jenkins X](https://jenkins-x.io)
- [Zeebe](https://helm.zeebe.io)


# Installation phase

First, we will install Jenkins X inside our Kubernetes Cluster. We will use this cluster to build and releasing and testing our applications. 
Jenkins X is divided into a CLI tool `jx` and a set of server-side components. 

We begin by installing `jx`, if you are running in Mac OSX you can use `brew`

```
```

Once we have `jx` installed we can run `jx boot` to install the server-side components into the Kubernetes cluster. Remember that `jx` will use the currently configured `kubectl` context, so make sure that you are pointing to the right cluster before running `jx boot`. 

Follow the steps proposed by `jx boot`, for reference these are the options that I've selected in a GKE cluster.  


# Scenario

During this workshop we will be helping a company that is in charge of helping conference organizers to setup their conference sites and back-office tools to setup awesome conferences. 

Their current application is a Monolith and it looks like this: 

The source code for this application can be [found here]()

The workshop aims to provide the tools, steps, and practices that can facilitate the migration from this application to a Cloud-Native architecture that runs on Kubernetes. In that Journey, we can enable teams to work independently and release their software in small increments while applying some of the principles outlined by the [Accelerate book](). 

This workshop covers the following steps:
- Running our Monolith in Kubernetes / CI/CD for our Monolith
- Introducing an API Gateway to serve as single entry points for all our services
- Split our Monolith into a set of Microservices
- Testing & Versioning a Cloud-Native Application
- Exposing business metrics and insights


# CI/CD for our Monolith

When moving to Kubernetes it is quite common to lift and shift our monolith applications that are already running outside Kubernertes. This will require to containarize our application and then provide all Kubernetes Manifests to deploy our application into an existing cluster. 

This exercise pushes us to learn Kubernetes basics concepts such as Deployments, Services, and Ingresses, as well as Docker basics such as how to build and publish a Docker Image and which base Docker Image should we use for our applications. While this is needed for deploying our applications into a running cluster, once we have done these steps for a couple of services/applications, we don't want to do them for 100 services. This is where Jenkins X comes to help us. 

You can find our [monolith application here]()






