# Crossplane and VCluster

Provisioning VClusters using Crossplane.

This tutorial shows how you can provision new VClusters using a Crossplane Composition. The idea here is to show how you can, with a single Kubernetes Cluster, provide isolation to different tenants while at the same time rely on the power of Crossplane Compositions to create Virtual Clsuters in the same way as you would create any other Cloud Resource. 

The advantages of using VCluster in contrast with creating a full-blown Kubernetes cluster is that we save on costs while we provide more isolation between resources and API Server calls than when we use `namespaces`. 


## Installation on KinD

```
kind create cluster
```

```
kubectl create ns crossplane-system
helm install crossplane --namespace crossplane-system crossplane-stable/crossplane
```


```
kubectl crossplane install provider crossplane/provider-helm:v0.10.0
```

```
SA=$(kubectl -n crossplane-system get sa -o name | grep provider-helm | sed -e 's|serviceaccount\/|crossplane-system:|g')
kubectl create clusterrolebinding provider-helm-admin-binding --clusterrole cluster-admin --serviceaccount="${SA}"
```

```
kubectl apply -f helm-provider-config.yaml
```

You need to install the `vcluster` CLI to connect to the cluster

[https://www.vcluster.com/docs/getting-started/setup](https://www.vcluster.com/docs/getting-started/setup)

## Creating VClusters using Crossplane Composition

```
kubectl apply -f composition/composition.yaml
kubectl apply -f composition/environment-resource-definition.yaml
```

```
kubectl apply -f composition/environment-resource.yaml
```

```
vcluster list 
```

```
vcluster connect dev-environment -- bash
```
or

```
vcluster connect dev-environment -- ksh
```

```
kubectl get pods -n conference
```

