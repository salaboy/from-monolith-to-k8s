# Creating your own platform

Create a Platform Cluster.

Then install Kratix: 

```
kubectl apply -f distribution/kratix.yaml
kubectl apply -f hack/platform/minio-install.yaml
```

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


Check: 
```
mc ls local/kratix-crds
```


## Reference
https://github.com/syntasso/kratix/blob/main/docs/quick-start.md
