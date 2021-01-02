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
kubectl apply -f https://github.com/salaboy/from-monolith-to-k8s/keycloak/k8s/keycloak-deployment.yaml -n sso

kubectl apply -f https://github.com/salaboy/from-monolith-to-k8s/keycloak/k8s/keycloak-service.yaml -n sso
```

Let's see the keycloak pod

```
$ kubectl get pods -n sso
```

## Changing API Gateway to secure our hidden microservices

[API Gateway](https://github.com/mcruzdev/fmtok8s-api-gateway) was created with Spring Cloud Gateway. The Spring Cloud Gateway uses Spring Webflux working with reactive stack.

There is a great lib called
`org.keycloak:keycloak-spring-boot-starter` that help us to configure our application using keycloak and it runs better with Servlet applications. [See](https://keycloak.discourse.group/t/webflux-support-for-spring-boot-and-spring-security-adapters/2936)

In this workshop, you will use Spring Security OAuth2. Let's go to use it.

### Adding Spring OAuth2 dependecies in API Gateway

```
<dependency>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-starter-oauth2-client</artifactId>
</dependency>
<dependency>
    <groupId>org.springframework.cloud</groupId>
    <artifactId>spring-cloud-starter-security</artifactId>
</dependency>
```

### Configuring 