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

This allows you to route traffic from http://localhost to services running inside the cluster. Notice that for KinD to work in this way, when we created the cluster we provided extra parameters and labels for the control plane node:
```
nodes:
- role: control-plane
  kubeadmConfigPatches:
  - |
    kind: InitConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "ingress-ready=true" #This allow the ingress controller to be installed in the control plane node
  extraPortMappings:
  - containerPort: 80 # This allows us to bind port 80 in local host to the ingress controller, so it can route traffic to services running inside the cluster.
    hostPort: 80
    protocol: TCP
  - containerPort: 443
    hostPort: 443
    protocol: TCP
```


## Installing the Application using Helm
Finally, we can install the application by adding a Helm chart repository. To achieve this, first we need to add a custom Helm Chart Repository for this application: 

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
LAST DEPLOYED: Thu Apr 21 09:04:46 2022
NAMESPACE: default
STATUS: deployed
REVISION: 1
TEST SUITE: None
NOTES:
Cloud-Native Conference Platform V1

Chart Deployed: fmtok8s-app - 0.1.0
Release Name: app

```

Now you can check that the pods of the application are being created correctly with: 
```
kubectl get pods
NAME                                       READY   STATUS      RESTARTS   AGE
app-fmtok8s-agenda-rest-65d7778fd7-pp4tf   1/1     Running     0          7m22s
app-fmtok8s-api-gateway-5dfd76594-j5bms    1/1     Running     0          7m22s
app-fmtok8s-c4p-rest-856c5d584d-9m7wl      1/1     Running     0          7m22s
app-fmtok8s-email-rest-7c8f54f6d9-s5ztz    1/1     Running     0          7m22s
```

When yo0u  get all the pods up and running, you should be able to access the application pointing your browser to: http://localhost


## Other options

If you don't want to create a Helm release, which gets created when you run `helm install` you can use Helm to produce the all YAML files of your application's services by running `helm template`. 

