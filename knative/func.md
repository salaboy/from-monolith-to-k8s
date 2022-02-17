# Knative `func` - Serverless experience on top of Kubernetes

This tutorial shows how to use the `func` CLI in conjunction with Knative Serving and Eventing to build applications using a serverless approach. 
Removing the need from developers to worry about Docker Containers or Kubernetes itself. 

## Pre requisites
- Kubernetes Cluster
- Knative Serving & Eventing installed
- `func` CLI installed
- (Optional) for OnCluster builds Tekton installed inside the cluster

## Scenario
This examples builds up on top of the [CloudEvents example](https://github.com/salaboy/from-monolith-to-k8s/tree/master/cloudevents), we will be creating two functions in different languages Java and Go. We will be using Knative Eventing to route events between them by defining Knative Brokers and Triggers, as we also did in the [CloudEvents example with Knative](https://github.com/salaboy/from-monolith-to-k8s/tree/master/cloudevents#with-knative-eventing). 


### Creating Java Function using Spring Boot + Spring Native

```
func create fmtok8s-java-function -l springboot -t cloudevents
```

This command creates a new directory called `fmtok8s-java-function` which contains a Spring Boot project using Spring Cloud Functions. It also uses Spring Native to generate a native binary that leverages GraalVM to enable Java applications to boot up faster to have shorter cold-starts. The `-t cloudevents` parameter specifies that we will use a template that contains a function that is ready to process a CloudEvent. The other alternative is to choose the `http` template which exposes a REST endpoint and can consume an HTTP request without requiring a CloudEvent specifically. 

You can use your favourite IDE to open the project and check the simple function that is generated from this template and edit it to your needs. 

If you take a look at the source code generated, right from the start we can see that the programming model is different in here. 

Spring Cloud Functions 

But if you want to get things going with the example function that is generated for you, you can run `func build` to build the function and generate a container for it, that by default will be pushed to the default Docker registry which is Docker Hub. 

```
salaboy> func build
A registry for Function images is required. For example, 'docker.io/tigerteam'.
? Registry for Function images: salaboy
Note: building a Function the first time will take longer than subsequent builds
🕕 Still building
   🙌 Function image built: docker.io/salaboy/fmtok8s-java-function:latest
```

And just like that, without worrying about having a Dockerfile, `func` uses [CNCF Buildpacks](http://buildpacks.io) to build and containarize your application. As you can see, I've provided my Docker Hub user name (`salaboy`) to automatically push the container to the registry, so then the container can be fetched from inside the cluster. 

The next step is to just run this function inside a configured cluster. Once again, `func` will take care creating the correct Knative Serving Service for our function to run.
```
salaboy> func deploy 
   🙌 Function image built: docker.io/salaboy/fmtok8s-java-function:latest
   Function deployed at URL: http://fmtok8s-java-function.default.X.X.X.X.sslip.io
```

There you go, the function is deployed and ready to accept requests at the following URL: `http://fmtok8s-java-function.default.X.X.X.X.sslip.io`
You can test your function by sending a CloudEvent using `curl` or `func invoke`.

```
curl -X POST http://fmtok8s-java-function.default.X.X.X.X.sslip.io
```

```
func invoke
```

## Creating a Go Function 

```
```

# Routing Events to functions using Knative Eventing Brokers and Triggers

```
kubectl create -f - <<EOF
apiVersion: eventing.knative.dev/v1
kind: Trigger
metadata:
  name: java-function-trigger
  namespace: default
spec:
  broker: default
  filter:
    attributes:
      type: uppercase
  subscriber:
    ref:
      apiVersion: serving.knative.dev/v1
      kind: Service
      name: fmtok8s-java-function
--- 

apiVersion: eventing.knative.dev/v1
kind: Trigger
metadata:
  name: go-function-trigger
  namespace: default
spec:
  broker: default
  filter:
    attributes:
      type: uppercase
  subscriber:
    ref:
      apiVersion: serving.knative.dev/v1
      kind: Service
      name: fmtok8s-go-function

EOF
```
