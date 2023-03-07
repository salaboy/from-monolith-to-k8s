# Code your CI Pipelines and Run them everywhere with Dagger

In this short tutorial we will be looking at the Dagger Service Pipelines included with each applicaation service. 
These pipelines are implemented in Go using the Dagger Go SDK and take care of building each service, creating a container, publishing it and creating a Helm chart that can be distributed for other teams to use. 

## Requirements

For running these pipelines locally you will need: 
- [Go installed](https://go.dev/doc/install)
- [A container runtime (such as Docker running locally)](https://docs.docker.com/get-docker/)

To run the pipelines remotely on a Kubernetes Cluster you can use [KinD](https://kind.sigs.k8s.io/) or any Kubernetes Cluster that you have available. 

## Let's run some pipelines

You can find each pipeline at each of the services repository's in a file called `pipeline.go`: 

- [Agenda Service](https://github.com/salaboy/fmtok8s-agenda-service/blob/main/pipelene.go)
- [Email Service](https://github.com/salaboy/fmtok8s-email-service/blob/main/pipeline.go)
- [Call for Proposals Service](https://github.com/salaboy/fmtok8s-c4p-service/blob/main/pipeline.go)
- [User Interface](https://github.com/salaboy/fmtok8s-frontend/blob/main/pipelene.go)

Feel free to clone any of these repositories locally by running for example: 

```
git clone https://github.com/salaboy/fmtok8s-email-service.git
cd fmtok8s-email-service/
```

You can run any defined task inside the `pipeline.go` file:

```
go mod tidy
go run pipeline.go <TASK>
```

The following tasks are defined for all the services: 
- `build` will build the service source code and create a container for it. This doesn't require any arguments. 
- `publish-image` publish the created container image to a container registry. This requires you to be logged in to the container registry and to provide the name of the container image as it will be published. For example: `salaboy/fmtok8s-agenda-service:0.1.0-dagger`. Notice that you need to include the org/username of where the image will be published.
- `helm-package` creates the Helm Chart ready to be distributed. 
- `publish-helm` uploads the Helm Chart to a Chart Repository. This require you to have the right credentials to connect and push to a Chart repository. 

If you run `go run pipeline.go full` all the tasks will be executed. Before being able to run all the tasks you will need to make sure that you have all the pre-requisites set, as for pushing to a Container Registry you will need to provide appropriate credentials. 

You can safely run `go run pipeline.go build` which doesn't require you to set any credentials. 

Now, for development purposes, this is quite convinient, because you can now build your service code in the same way that your CI (Continupis Integration) system will do. But you don't want to run in production container images that were created in your developer's laptop right? 
The next section shows a simple setup of running Dagger pipelines remotely inside a Kubernetes Cluster. 

## Running your pipelines remotely on Kubernetes

The Dagger Pipeline Engine can be run anywhere where you can run containers, that means that it can runs in Kubernetes without the need of complicated setups. 
In this short tutorial we will run the pipelines that we were running locally with our local container runtime, now remotely against a Dagger Pipeline Engine that runs inside a Kuberneets Pod. This is an experimental feature, and not a recommended way to run Dagger, but it help us to prove the point. 

Let's run the Dagger Pipeline Engine inside Kubernetes by creating a Pod with Dagger: 

```
kubectl run dagger image=registry.dagger.io/engine:v0.3.13 --privileged=true
```

Alternatively, you can apply the `k8s/pod.yaml` manifest using `kubectl apply -f k8s/pod.yaml`.

**Note**: this is far from ideal because we are not setting any persistence or replication mechanism for Dagger itself, all the caching mechanism are volatile in this case. Check the official documentation for more about this. 

Now to run the projects pipelines against this remote service you only need to export the following environment variable: 
```
export _EXPERIMENTAL_DAGGER_RUNNER_HOST=kube-pod://<podname>?context=<context>&namespace=<namespace>&container=<container>
```

Where `<podname>` is `dagger` (because we created the pod manually), `<context>` is your Kubernetes Cluster context, if you are running against a KinD Cluster this might be `kind-kind`. You can find your current context name by running `kubectl config current-context`. Finally `<namespace>` is the namespace where you run the Dagger Container, and `<container>` is once again `dagger`. For my setup against KinD, this would look like this: 

```
export _EXPERIMENTAL_DAGGER_RUNNER_HOST="kube-pod://dagger?context=kind-kind&namespace=default&container=dagger"
```

Notice also that my KinD cluster didn't had anything related to Pipelines. 

Now if you run in any of the projects: 
```
go run pipeline.go build 
```

The build will happen remotely inside the Cluster. If you were running this against a remote Kubernetes Cluster (not KinD), there will not be need for you to have a local Container Runtime to build your services and their containers. 
