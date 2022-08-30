# Argo Rollouts for Release Strategies

In this short tutorial, we will take a look at Argo Rollouts and how they helps us to implement different release strategies. We will be just using a single service to demonstrate how the Argo Rollouts Rollouts resources works, but you can definitely add more services to the mix. 

We will be also creating a local KinD cluster where we will be installing Argo Rollouts and our target service. 

## Installation and getting started

Let's start by creating our KinD Cluster with Ingress Enabled (we need an ingress to test loadbalancing):

Create a KinD Cluster:

```
cat <<EOF | kind create cluster --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  kubeadmConfigPatches:
  - |
    kind: InitConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "ingress-ready=true"
  extraPortMappings:
  - containerPort: 80
    hostPort: 80
    protocol: TCP
  - containerPort: 443
    hostPort: 443
    protocol: TCP
EOF
```

We need NGINGX Ingress Controller to route traffic from our laptop to the services that are running inside the cluster. NGINX Ingress Controller act as a router that is running inside the cluster, but exposed to the outside world. 

```
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
```

This allows you to route traffic from http://localhost to services running inside the cluster.

Let's install Argo Rollouts now following the [guide in the official documentation](https://argoproj.github.io/argo-rollouts/installation/#controller-installation)

Also install the Argo Rollouts `kubectl` plugin. 

Locally, against the KinD cluster I've just executed these commands: 

```
kubectl create namespace argo-rollouts
kubectl apply -n argo-rollouts -f https://github.com/argoproj/argo-rollouts/releases/latest/download/install.yaml
```

And because I am running on Mac OSX I used `brew` to install the `kubectl` plugin: 

```
brew install argoproj/tap/kubectl-argo-rollouts
```

Once you have the plugin you can start a local version of the Argo Rollouts Dashboard, by running in a new terminal: 

```
kubectl argo rollouts dashboard
```

Then you can access the dashboard by pointing your browser to [http://localhost:3100](http://localhost:3100)



## Why should we use Argo Rollouts?

Kubernetes provides us the base building blocks to make sure that we can keep our applications running without any downtimes. But dealing with Deployments and Services to implement different release strategies is complicated and prone to error. For these reasons tools like Argo Rollouts were created. Argo Rollouts implement different high-level release strategies and the mechanisms to rollout our workloads following a progressive delivery approach. 

In the following sections we will see some of the Argo Rollouts capabilities to implement: 
- Canary Releases
- Blue-Green Deployments

## Canary Releases for our Email Service

When using Argo Rollouts, we delegate the deployment side of things to a new CRD called `Rollout`. This new CRD from Argo Rollouts will define the release strategy that we want to use and how it will perform the rollout of new versions of our service, in this case our Email Service.

Our Rollout resources is not the one responsible for defining what needs to be deployed, but also how it will be rolled out when new updates happen: 


```
kubectl apply -f canary-release/
```

Watch rollout status

```
kubectl argo rollouts get rollout email-service-canary --watch
```

@TODO: do rollout edit to change the version too
```
kubectl argo rollouts set image email-service-canary \
  email-service=ghcr.io/salaboy/fmtok8s-email-service:v0.2.0-native
```


```
kubectl argo rollouts promote email-service-canary
```



```
kubectl argo rollouts abort email-service-canary
```

# Blue/Green Deployments

```
kubectl apply -f blue-green/
```


Watch rollout status

```
kubectl argo rollouts get rollout email-service-bluegreen --watch
```

@TODO: do rollout edit to change the version too
```
kubectl argo rollouts set image email-service-bluegreen \
  email-service=ghcr.io/salaboy/fmtok8s-email-service:v0.2.0-native
```

Promote the rollout

```
kubectl argo rollouts promote email-service-bluegreen
```
