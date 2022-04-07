# Creating your own platform


## Create a Platform Cluster:

Then install Kratix: 

```
kubectl apply -f distribution/kratix.yaml
kubectl apply -f hack/platform/minio-install.yaml
```

**Note**: change minio service to LoadBalancer type to expose outside of GKE

We are using Minio inside the Cluster to store resources, 
Port Forward
```
k port-forward svc/minio -n kratix-platform-system 9000:80
```
Then configure local minio-mc client to connect: 
```
vi /Users/msalatino/.mc/config.json

```
Configure local keys to `minioamdin`/`minioamdin`


Install a promise like Jenkins Promise: 

```
kubectl apply -f samples/jenkins/jenkins-promise.yaml
```


Check: 
```
mc ls local/kratix-crds
```

## Create a Worker Cluster



## Reference
https://github.com/syntasso/kratix/blob/main/docs/quick-start.md
