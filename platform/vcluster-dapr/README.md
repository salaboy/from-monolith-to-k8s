# VCluster and Dapr


This short tutorial shows how we can reuse a single Dapr installation across multiple VClusters. 
This uses the `multiNamespaceMode` and the `syncer` extension mechanism to import/export Dapr components from and into VClusters. 

## Prerequisites

- Create a KinD Cluster
- VCluster CLI
- Install Dapr
- Install Redis on Host Cluster

Create a KinD Cluster:

```
kind create cluster
```

Install Dapr: 


```
helm repo add dapr https://dapr.github.io/helm-charts/
helm repo update
helm upgrade --install dapr dapr/dapr \
--version=1.10.4 \
--namespace dapr-system \
--create-namespace \
--wait
```

Once Dapr is installed we need to customize the Sidecar injector deployment located in the `dapr-system` namespace: 

```
kubectl edit deploy -n dapr-system dapr-sidecar-injector
```

Make the following changes: 

```
image: image: daprio/injector:nightly-2023-03-15-linux-amd64
```

and: 

```
env: 
- name: ALLOWED_SERVICE_ACCOUNTS_PREFIX_NAMES
  value: vcluster-dapr-enabled:vc-dapr-enabled
```

**Notice that `dapr-enabled` is the name of the vcluster that we will be creating**. 


Let's install Redis to be used as our Statestore implementation:

```
helm install redis bitnami/redis --set architecture=standalone
```

Get the redis password: 
```
export REDIS_PASSWORD=$(kubectl get secret --namespace default redis -o jsonpath="{.data.redis-password}" | base64 -d) 
echo $REDIS_PASSWORD
```

You need to add this password to the `statestore.yaml` file. 

## Let's create some Dapr-enabled VClusters

Let's create a new VCluster using the configuration located in the `values.yaml` file: 

```
vcluster create dapr-enabled -f values.yaml
```

You are now inside your VCluster, let's create a statestore: 

```
kubectl apply -f statestore.yaml
```

Let's deploy a couple of applications using the Statestore: 

```
kubectl apply -f apps.yaml
```

This should start two applications that connect to the statestore component that we created inside the VCluster, but connects to the Host Redis instance. 

