# Creating your own platform

In this document we will be using Kratix, Crossplane, Knative and RabbitMQ to create a custom platform configuration. 



## Create a Platform Cluster:

Then install Kratix: 

```
kubectl apply -f distribution/kratix.yaml
kubectl apply -f hack/platform/minio-install.yaml
```


We are using Minio inside the Cluster to store resources, 
If you want to access from your local environment then use `Port Forward`

```
k port-forward svc/minio -n kratix-platform-system 9000:80
```

**Note**: To expose minio service to the worker clusters change the `serviceType` to `LoadBalancer` to expose outside of GKE. Check the Service IP and use it to configure the Worker Cluster resources. 


Then configure local minio-mc client to connect: 
```
vi /Users/msalatino/.mc/config.json

```

Configure local keys to `minioamdin`/`minioamdin`


Install a promise like Jenkins Promise: 

```
kubectl apply -f samples/jenkins/jenkins-promise.yaml
```


Check Minio, this creates the CRDs inside Minio for the Worker Cluster (using flux) to pick up and apply: 

```
mc ls local/kratix-crds
```

## Create a Worker Cluster

GKE Create Cluster, this will install flux and the GitOps tools:

```
kubectl apply -f hack/worker/gitops-tk-install.yaml
kubectl apply -f hack/worker/gitops-tk-resources.yaml
```

For this to work, in the `gitops-tk-resources.yaml` we need to change the minio endpoints IP to the Minio Service public IP in [here](https://github.com/syntasso/kratix/blob/main/hack/worker/gitops-tk-resources.yaml#L11) and [here](https://github.com/syntasso/kratix/blob/main/hack/worker/gitops-tk-resources.yaml#L25). This will enable flux to pick up the minio resources to sync. 






## Reference
https://github.com/syntasso/kratix/blob/main/docs/quick-start.md
