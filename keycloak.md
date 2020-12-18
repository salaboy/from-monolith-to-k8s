# From Monolith to K8s - Workshop 


# Installing and Configuring Keycloak

During this step-by-step you will be using **Kubernetes Cluster** and a Keycloak as SSO to secure our API Gateway and Microservices. 

### Creating a Kubernetes Cluster with KIND

```
$ kind create cluster --name keycloak
```

Don't forget to set current cluster/context

```
$ kubectl cluster-info --context kind-keycloak
```

We will create a namespace sso

```
$ kubectl create namespace sso
```

### Adding Keycloak on Cluster Kubernetes

```
kubectl apply -f https://github.com/salaboy/from-monolith-to-k8s/keycloak/k8s
```

