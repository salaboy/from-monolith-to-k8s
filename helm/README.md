# Helm as a Package Manager

You can use [Helm](http://helm.sh) to package, distribute and install your Kubernetes Services or applications. 

This document shows how easy is to install an application that has been packaged as a helm chart. 

## Installing Helm

Helm is a client-side only tool that will connect to your Kubernetes Clusters using the `./kube/config` access tokens. 

You can install Helm by following the [instructions here](https://helm.sh/docs/intro/install/).

Once you have Helm installed, you can add **Helm Repositories** which contains Helm Charts (packages).

For this tutorial you will be adding my Helm Repository which is hosted in a GitHub Site: [https://salaboy.github.io/helm/](https://salaboy.github.io/helm/)

You can add this Helm repository to you Helm installation by running:

```
salaboy> helm repo add fmtok8s https://salaboy.github.io/helm/
```

You will need to run `helm update` everytime you want to fetch the latest version of the charts or the first time after adding a new repository: 

```
salaboy> helm repo update
```

## Installing the Conference Application

You can always do searchs against the added Helm repositories with `helm search repo <Name>`

For example, you can search the following application, which is hosted in my repository `fmtok8s-app`: 

```
salaboy> helm search repo fmtok8s-app
```

Should return: 

```
NAME               	CHART VERSION	APP VERSION	DESCRIPTION                               
fmtok8s/fmtok8s-app	0.1.0        	0.1.0      	A Helm chart for a Conference Platform App
```

Which now you can install by running: 

```
salaboy> helm install app fmtok8s/fmtok8s-app

```

You should see something like: 
```
NAME: app
LAST DEPLOYED: Sat Jul  3 14:27:02 2021
NAMESPACE: default
STATUS: deployed
REVISION: 1
TEST SUITE: None
NOTES:
Cloud-Native Conference Platform V1

Chart Deployed: fmtok8s-app - 0.1.0
Release Name: app

```

The application is composed by 4 independent services that can be installed separately, but the `fmtok8s-app` chart install them all configured to work together. The `fmtok8s-app` chart also comes with an `Ingress` defintion called `frontend`. For this to work you will need to have an Ingress Controller deployed in the cluster. 

You can check that the application is up and running by checking if the Kubernetes Pods have started and are in READY status: 

```
salaboy> kubectl get pods
```

You should see something like: 
```
NAME                                       READY   STATUS      RESTARTS   AGE
app-fmtok8s-agenda-rest-68bd9c8bcb-8njzf   1/1     Running     0          2m18s
app-fmtok8s-api-gateway-58d49588b4-4psdh   1/1     Running     0          2m18s
app-fmtok8s-c4p-rest-7cb8bc4485-2l5h2      1/1     Running     0          2m18s
app-fmtok8s-email-rest-8f954fbbd-99nkb     1/1     Running     0          2m18s
```



