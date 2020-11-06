# Workshop 

During this workshop you will deploy a Cloud Native application, inspect it, change its configuration to use different services and 
play around with it to get familiar with Kubernetes and Cloud Native tools that can help you to be more efficient. 

During this workshop you will be using GKE (Managed Kubernetes Engine inside Google Cloud) to deploy a complex application composed by multiple services. 

@TODO: Add GKE instructions here

In order to do this you will be using `kubectl` and `helm` to deploy the application. Because you will be using the Google Cloud Console, you can save some time by creating some aliases for these two commands:

```
> alias k=kubectl
> alias h=helm
```
Now instead of typing `kubectl` or `helm` you can just type `k` and `h` respectivily. 

# Pre Requisites
- Knative Service

```
> k apply --filename https://github.com/knative/serving/releases/download/v0.18.0/serving-crds.yaml
> k apply --filename https://github.com/knative/serving/releases/download/v0.18.0/serving-core.yaml
> k apply --filename https://github.com/knative/net-kourier/releases/download/v0.18.0/kourier.yaml
> k patch configmap/config-network \
  --namespace knative-serving \
  --type merge \
  --patch '{"data":{"ingress.class":"kourier.ingress.networking.knative.dev"}}'
> k apply --filename https://github.com/knative/serving/releases/download/v0.18.0/serving-default-domain.yaml
```

- Knative Eventing

```
> k apply --filename https://github.com/knative/eventing/releases/download/v0.18.0/eventing-crds.yaml
> k apply --filename https://github.com/knative/eventing/releases/download/v0.18.0/eventing-core.yaml
> k apply --filename https://github.com/knative/eventing/releases/download/v0.18.0/in-memory-channel.yaml
> k apply --filename https://github.com/knative/eventing/releases/download/v0.18.0/mt-channel-broker.yaml
> k create -f - <<EOF
apiVersion: eventing.knative.dev/v1
kind: Broker
metadata:
 name: default
 namespace: default
EOF
```

- Create a Camunda Cloud account, you will need it for version 2 and version 3 of the applications. Sign up for [Camunda Cloud](https://accounts.cloud.ultrawombat.com/signup?campaign=workshop)

# Installing our Cloud Native Application

Once you have the aliases, Knative and a Camunda Cloud account  you can proceed to add a new Helm Repository where the Helm packages for the application are stored. 
You can do this by runnig the following command: 

```
> h repo add workshop http://chartmuseum-jx.35.222.17.41.nip.io
> h repo update
```

Now you are ready to install the application by just running the following command:
```
> h install fmtok8s workshop/fmtok8s-app
```

# Understanding our application
@TODO: explain from a high level what the first version of the application is doing and some of the challenges. 


## Deploying Version 2

Go to the Camunda Cloud Console, create a cluster and a client. Copy the credentials Kubernetes Secret command from the client popup and paste it into the Google Cloud Console: 
```
k create secret generic camunda-cloud-secret --from-literal=ZEEBE_ADDRESS=<ZEEBE_ADDRESS HERE> --from-literal=ZEEBE_CLIENT_ID=<ZEEBE_CLIENT_ID HERE> --from-literal=ZEEBE_CLIENT_SECRET=<ZEEBE_CLIENT_SECRET HERE> --fr
om-literal=ZEEBE_AUTHORIZATION_SERVER_URL=<ZEEBE_AUTHORIZATION_SERVER_URL HERE>
```

@TODO: explain that V2 of the application is emitting events using Knative, you can observe these to understand what is happening in your application.
@TODO: explain Zeebe Cloud Events Router

```
> h install fmtok8s workshop/fmtok8s-app-v2
```

# Workflow Orchestration with Camunda Cloud

```
> h install fmtok8s workshop/fmtok8s-app-v3
```




