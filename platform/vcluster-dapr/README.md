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
image: daprio/injector:nightly-2023-03-15-linux-amd64
```

and: 


```
env: 
- name: ALLOWED_SERVICE_ACCOUNTS_PREFIX_NAMES
  value: vcluster-*:vc-*
```

Let's install Redis to be used as our Statestore implementation:

```
helm install redis bitnami/redis --set architecture=standalone
```

Get the redis password: 
```
export REDIS_PASSWORD=$(kubectl get secret --namespace default redis -o jsonpath="{.data.redis-password}" | base64 -d) 
echo $REDIS_PASSWORD
```

You need to update this password in the `statestore.yaml` file. 

## Let's create some Dapr-enabled VClusters

Let's create a new VCluster using the configuration located in the `values.yaml` file: 

```
vcluster create --chart-version 0.15.0-alpha.0 dapr-enabled -f values.yaml
```

Once the VCluster is created, we are automatically connected to it. 

You can check this by running:
```
kubectl get ns
```
You should see, something like this: 
```
> kubectl get ns
NAME              STATUS   AGE
default           Active   48s
kube-system       Active   48s
kube-public       Active   48s
kube-node-lease   Active   48s
```

Notice that there is no `dapr-system` namespace, as we installed Dapr in the Host cluster. 

You are now inside your VCluster, let's create a statestore: 

```
kubectl apply -f statestore.yaml
```

Let's deploy a couple of applications using the Statestore: 

```
kubectl apply -f apps.yaml
```

This should start two applications that connect to the statestore component that we created inside the VCluster, but connects to the Redis instance installed in the Host cluster.


To interact with these two applications you will need to use `port-forward` to both services.

To access the Java Application using port 8080:
```
kubectl port-forward svc/java-app-service 8080:80 
```
To access the Go Application using port 8081:
```
kubectl port-forward svc/go-app-service 8081:80
```

Now you can send request to both services. If you send first a request to the Go service that reads data from Redis you will see an empty array back. 

```
> curl localhost:8081          
{"Values":null} 
```

Let's now send a request to the Java application to store a value: 

```
> curl -X POST "localhost:8080/?value=42"
{"values":["42"]}
```

You can try adding another value. Notice that the request returns the currently stored values too.

Then the Go Application can return the stored value: 

```
> curl localhost:8081
{"Values":["42"]}
```

If you want to delete all values you can send a `DELETE` request to the Java Application: 

```
> curl -X DELETE localhost:8080
```
