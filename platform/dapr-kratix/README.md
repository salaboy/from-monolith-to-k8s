# Building Dapr-enabled Platforms using Kratix

In this tutorial we will use Kratix to configure our Dapr enabled development environment. 
The main objective is to request new development environments for our applications that comes with all the components needed already pre-installed and configured to work. 

## Installing Kratix 

We will be using a Kubernetes KinD Cluster to run Kratix and configure our development environments. You can create a KinD Cluster by running: 

```
kind create cluster --name platform --image kindest/node:v1.24.7
```

Then we install Kratix into the cluster by following their getting started guide: https://kratix.io/docs/main/quick-start

```
kubectl apply --filename https://raw.githubusercontent.com/syntasso/kratix/main/distribution/single-cluster/install-all-in-one.yaml
```

And then registering the same cluster as a worker cluster. If you are working in multiple clusters, this  

```
kubectl apply --filename https://raw.githubusercontent.com/syntasso/kratix/main/distribution/single-cluster/config-all-in-one.yaml
```
# Defining Dapr-enabled Development environments

We will create a new Kratix promise to represent the components that needs to be installed for each Development Environment that we want to create. 
For this example, we want a developer environment to have Dapr installed, Redis installed and a Dapr Statestore component ready to be used by our developers. 

For that we create a new Kratix promise that installs the Dapr and Redis promise and defines our Dapr Statestore Component configuration to point to the installed Redis instance.

```env-promise/promise.yaml
apiVersion: platform.kratix.io/v1alpha1
kind: Promise
metadata:
  name: env
  namespace: default
spec:
  workerClusterResources:
  - apiVersion: platform.kratix.io/v1alpha1
    kind: Promise
    metadata:
      creationTimestamp: null
      name: dapr
      namespace: default
    ...
  - apiVersion: platform.kratix.io/v1alpha1
    kind: Promise
    metadata:
      creationTimestamp: null
      name: redis
      namespace: default
    ...
  xaasCrd:
    apiVersion: apiextensions.k8s.io/v1
    kind: CustomResourceDefinition
    metadata:
      name: env.marketplace.kratix.io
    spec:
      group: marketplace.kratix.io
      names:
        kind: env
        plural: env
        singular: env
      scope: Namespaced
      versions:
      - name: v1alpha1
        schema:
          openAPIV3Schema:
            properties:
              spec:
                properties:
                  deployApp:
                    description: |
                      Deploy the Dapr Applications defined
                    type: boolean
                  database:
                    properties: 
                      statestoreName:
                        description: |
                          The name of the statestore component to create
                        type: string
                      enabled:
                        description: |
                          Deploy an instance of Redis or not
                        type: boolean
                    type: object
                type: object
            type: object
        served: true
        storage: true
  xaasRequestPipeline:
  - salaboy/environment-request-pipeline:v0.1.0

```

Then we apply this promise to the cluster.

```bash
kubectl apply -f env-promise/promise.yaml
```

This will install Dapr and Redis Promises, which themselves installed Dapr and
the Redis operator.


# Requesting new Dapr-enabled Development Environment

Now that we have our promise installed in Kratix we can create new development environments by sending the following resource:

```env-promise/my-dev-env.yaml
apiVersion: marketplace.kratix.io/v1alpha1
kind: env
metadata:
  name: my-dev-env
  namespace: default
spec:
  deployApp: true
  database:   
    enabled: true
    statestoreName: "statestore"

```

Then we then make a new Development Environment request to the Platform cluster:

```bash
kubectl apply -f env-promise/my-dev-env.yaml
```

Now if we inspect the cluster we will see we have a Redis instance running, and
the components configurated in Dapr:
```
kubectl get pods
```

```
kubectl get components
```

You should see something like this: 

```
```

# Interacting with the deployed applications

Now that the applications are running we can use `kubect port-forward` to interact with both applications. 

In a separate terminal run: 
```
kubectl port-forward svc/java-app 8080:80
```

In another terminal run: 

```
kubectl port-forward svc/go-app 8081:80
```

Let's send an HTTP POST request to the Java App that writes to the Dapr Statestore component which is backed up by the Redis instance that the Kratix promise created and configured.

```
curl -X POST localhost:8080/?value=42
```

You should see something like this: 
```
```

Let's now send a GET request to the Go App which reads from the same Dapr Statestore component: 

```
curl localhost:8081
```

You should see something like this: 

```
```

You can find the source code of this application here: 
- Java App: check [here]() to see how the Dapr Java SDK is being used to connect to the Statestore by just using the statestore name. Notice that the name used is the same as the name of the Dapr component listed by running `kubectl get components`
- Go App: check [here]() to see how the Dapr Go SDK is being used to connect to the Dapr Statestore component.

The same principles can be used to expand this demo to use the Pub/Sub component to exchange message between applciations written in different languages, without pushing the applications to know which backing implementation is being used. This allow the application code to work across different implementaitons and even across cloud providers.


# Developing from Source

If you are working with this example, and you want to change the Kratix Pipeline you can find it in the `internal/request-pipeline` directory. You can build the pipeline into a container image by running, from within the `internal/scripts` directory: 

```
./pipeline-image build push
```

You can push your own version under your own registry, but remember to change this into the `promise.yaml` file, which makes reference to the pipeline container image that is going to be used by the Kratix Promise. 

# Cleaning up

You can delete your KinD Cluster to clean up all the installed components: 
```
kind delete clusters platform
```