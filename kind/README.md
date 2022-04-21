# Creating a KinD (Kubernetes in Docker) Cluster and installing the app

In this short tutorial we will be installing the Conference Application using Helm into a Kubernetes Cluster provisioned using KinD. This Kubernetes Cluster will run in our local machine, using our Docker Deamon and its configurations (CPUs and RAM allocations). 

## Pre Requisites
- Docker (check docker configurations for CPU and RAM allowences) 
- KinD
- Helm

## Creating our KinD Cluster

Create a KinD Cluster with 3 worker nodes and 1 Control Plane

```
cat <<EOF | kind create cluster --name dev --config=-
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
- role: worker
- role: worker
- role: worker
EOF

```

## Installing NGINX Ingress Controller
We need NGINGX Ingress Controller to route traffic from our laptop to the services that are running inside the cluster. NGINX Ingress Controller act as a router that is running inside the cluster, but exposed to the outside world. 

```
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/static/provider/kind/deploy.yaml
```

## Install Application using Helm
Finally, we can install the application by adding a Helm chart repository. 

Add a custom Helm Chart Repository for this application: 

```
helm repo add fmtok8s https://salaboy.github.io/helm/
helm repo update
```

Then run `helm install`: 

```
helm install app fmtok8s/fmtok8s-app
```

You should see the following output: 

```
NAME: app
LAST DEPLOYED: Sat Aug 28 13:52:48 2021
NAMESPACE: default
STATUS: deployed
REVISION: 1
TEST SUITE: None
NOTES:
Cloud-Native Conference Platform V1

Chart Deployed: fmtok8s-app - 0.1.0
Release Name: app

```

You should be able to access the application pointing your browser to: http://localhost:8080

