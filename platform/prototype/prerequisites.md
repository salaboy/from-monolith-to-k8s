# Prerequisites and Installation 

This tutorial creates and interacts with Kubernetes clusters, as well as installs Helm Charts. Hence, the following tools are needed: 
- [Install `kubectl`](https://kubernetes.io/docs/tasks/tools/)
- [Install `kind`](https://kind.sigs.k8s.io/docs/user/quick-start/#installation)
- [Install `helm`](https://helm.sh/docs/intro/install/) 
- [Install `docker`](https://docs.docker.com/engine/install/)

For this demo, we will create a Kubernetes cluster to host the platform tools which  create development environments. These tools use VCluster to create development environments in separate namespaces. We will also create one namespace for our Production Environment. For simplicity, this tutorial uses [KinD](https://kind.sigs.k8s.io/), but we encourage you to try the tutorial on a real Kubernetes cluster. 

- [Installing Command-Line Tools](installing-clis.md)
- [Create a Platform Cluster & Install Tools](platform-cluster.md)
  

## Creating and Configuring Our Platform Cluster

For this demo, our platform will enable development teams to request new `Environment`s by interacting with our Platform APIs. 

To host all the Platform tools we will use a KinD cluster (local). Inside this cluster we will install the following tools: 
- [Crossplane](https://crossplane.io)
- [Dapr](https://dapr.io)
- [ArgoCD](https://argoproj.github.io/)

For the sake of simplicity, we will be also creating a `production` namespace to simulate a different environment. 

You can create a KinD cluster install all these tools by following the steps here: [Platform Cluster installation](platform-cluster.md).

Once we have installed all the platform tools, we can configure them to work for our teams, and in this particular case we will optimize for scenarios focused on functions.

To enable teams to request `Environment`s in a self-service approach we have configured Crossplane in conjuction with VCluster. These `Environment`s can each be configured differently depending what the team needs to do. For this demo, I have created a [Crossplane Composition](https://crossplane.io/docs/v1.9/concepts/composition.html) that uses [VCluster](https://www.vcluster.com/) to create one virtual cluster per development environment requested. This enables a team to request their own isolated cluster so that they can work on features without clashing with other teams' work. 

For this to work, we need to create two things: the Custom Resource Definition (CRD) that defines the APIs for creating new `Environment`s, and the Crossplane Composition that defines the resources that will be created every time that a new `Environment` resource is created. 

Let's apply the Crossplane Composition and our **Environment Custom Resource Definition (CRD)** into the Platform Cluster:
```
kubectl apply -f crossplane/environment-resource-definition.yaml
kubectl apply -f crossplane/composition-devenv.yaml
```

The Crossplane Composition that we have defined and configured in our Platform Cluster uses the Crossplane Helm Provider to create a new VCluster for every `Environment` with `type: development`. The VCluster will be created inside the Platform Cluster, but it will provide its own isolated Kubernetes API Server for the team to interact with. 

The VCluster created for each development `Environment` is using the VCluster Knative Serving Plugin to enable teams to use Knative Serving inside the virtual cluster, but without having Knative Serving installed. The VCluster Knative Serving plugin shares the Knative Serving installation in the host cluster with all of the virtual clusters.

Now we are ready to request environments, deploy our applications/functions, and promote them to production. You can go back to the main tutorial to [request a new development environment](README.md#requesting-a-new-environment).