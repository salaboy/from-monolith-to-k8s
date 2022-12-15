# Using Dapr to build Cloud Native Apps

In this tutorial we will build stateful Functions that leverage the [dapr.io](https://dapr.io) stateful memory store. 

We will use Dapr in conjuction with [Knative Functions](https://knative.dev) to leverage scaling to zero and autoscaling provided by [Knative Serving](https://knative.dev).

## Pre-requisites & installation


We will be creating a local KinD cluster where we will install Knative Serving and Dapr.

For this you will need to install the following CLIs:

- [Install `kubectl`](https://kubernetes.io/docs/tasks/tools/)
- [Install `kind`](https://kind.sigs.k8s.io/docs/user/quick-start/#installation)
- [Install `helm`](https://helm.sh/docs/intro/install/) 
- [Install `docker`](https://docs.docker.com/engine/install/)
- [Install the Knative Functions `func` CLI](https://knative.dev/docs/functions/install-func/)
- [Install the `dapr` CLI](https://docs.dapr.io/getting-started/install-dapr-cli/)

Let's create a cluster: 

```
cat <<EOF | kind create cluster --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  extraPortMappings:
  - containerPort: 31080 # expose port 31380 of the node to port 80 on the host, later to be use by kourier or contour ingress
    listenAddress: 127.0.0.1
    hostPort: 80
EOF
```

Let's now install Knative Serving into the cluster: 

[Check this link for full instructions from the official docs](https://knative.dev/docs/install/yaml-install/serving/install-serving-with-yaml/#prerequisites)

```
kubectl apply -f https://github.com/knative/serving/releases/download/knative-v1.8.0/serving-crds.yaml
kubectl apply -f https://github.com/knative/serving/releases/download/knative-v1.8.0/serving-core.yaml

```

Installing the networking stack to support advanced traffic management: 

```
kubectl apply -f https://github.com/knative/net-kourier/releases/download/knative-v1.8.0/kourier.yaml

```

```
kubectl patch configmap/config-network \
  --namespace knative-serving \
  --type merge \
  --patch '{"data":{"ingress-class":"kourier.ingress.networking.knative.dev"}}'

```

Configuring domain mappings: 

```
kubectl apply -f https://github.com/knative/serving/releases/download/knative-v1.8.0/serving-default-domain.yaml

```

**Only for Knative on KinD** 

For Knative Magic DNS to work in KinD you need to patch the following ConfigMap:

```
kubectl patch configmap -n knative-serving config-domain -p "{\"data\": {\"127.0.0.1.sslip.io\": \"\"}}"
```

and if you installed the `kourier` networking layer you need to create an ingress:

```
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Service
metadata:
  name: kourier-ingress
  namespace: kourier-system
  labels:
    networking.knative.dev/ingress-provider: kourier
spec:
  type: NodePort
  selector:
    app: 3scale-kourier-gateway
  ports:
    - name: http2
      nodePort: 31080
      port: 80
      targetPort: 8080
EOF
```

Finally let's use install Dapr using Helm. Notice that you can also install Dapr using the `dapr` CLI.

```
helm repo add dapr https://dapr.github.io/helm-charts/
helm repo update
helm upgrade --install dapr dapr/dapr \
--version=1.9 \
--namespace dapr-system \
--create-namespace \
--wait
```

# Let's code!

We will be using some of the Dapr building blocks to enchance our Knative Functions and build an awesome application. 

![dapr-building-blocks.png](dapr-building-blocks.png)

We will start by creating a function using the Go programming language that store state into a Redis database by using the State Management building block. 

## Using the Dapr State Store 

To use the Dapr state store abstraction, first we need to define a Dapr component. For this example I've chosen Redis to be the implementation used to store state, but the idea here is that you can swap it to another state store if you need to without changing the code that is storing state.

```statestore.yaml
apiVersion: dapr.io/v1alpha1
kind: Component
metadata:
  name: statestore
spec:
  type: state.redis
  version: v1
  metadata:
  # These settings will work out of the box if you use `helm install
  # bitnami/redis`.  If you have your own setup, replace
  # `redis-master:6379` with your own Redis master address, and the
  # Redis password with your own Secret's name. For more information,
  # see https://docs.dapr.io/operations/components/component-secrets .
  - name: redisHost
    value: redis-master:6379
  - name: redisPassword
    secretKeyRef:
      name: redis
      key: redis-password
auth:
  secretStore: kubernetes
```

Before applying this resource to our cluster, we need to install Redis, and we will do this using Helm:

```
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
helm install redis bitnami/redis --set image.tag=6.2 --set architecture=standalone
```

Now we have a redis instance and a secret called `redis` which contains a key called `redis-password` that our `statestore` component will use to connect to it. Now you can go ahead and apply our `statestore` component into our cluster: 

```
kubectl apply -f statestore.yaml
```

## Let's create a function to save some state

```
func create -l go -t http
```


Once the function is created we can use the state store APIs to store state. These APIs are provided by Dapr and can be accessed using HTTP or GRPC. In this example we will use the HTTP endpoints that can be accessed at: 

```
http://localhost:3500/v1.0/state/statestore
```

From our function. No driver/library is needed to store state, just a simple HTTP request. As you might guess we will use a POST request to store state and GET requests to get the values that we stored.

In Go, this is as simple as sending a JSON payload using `http.Post()`

```functions/generate-values/handle.go
resp, err := http.Post(stateStoreUrl, "application/json",
		bytes.NewBuffer(json_data))
	if err != nil {
		log.Fatal(err)
	}
```

Notice also that `json_data` needs to be an array of objects, for example:  `[{ "key": "value"}]`.



In order to let Dapr know about our function, we need to make sure that the Knative Service generated by `func` include some Dapr annotations. For that reason we need to tweak the `func.yaml` file: 

```
...
deploy:
  namespace: default
  annotations:
    dapr.io/app-id: generate-values
    dapr.io/app-port: "8080"
    dapr.io/enable-api-logging: "true"
    dapr.io/enabled: "true"
    dapr.io/metrics-port: "9099"
...

```

Notice that I needed to change the `metrics-port` for Dapr, because the Knative `queue-proxy` uses the same port for metrics.

Once we have our request to store state, we can deploy our function to our cluster. 

```
func deploy -v -r docker.io/<username>
```

Now that the function is deployed you can interact with it by sending a `HTTP` request using `curl`:

```
curl http://generate-values.default.127.0.0.1.sslip.io
```

A random value was generated and stored into redis. 

First we need to obtain the REDIS_PASSWORD for the redis instance that we installed in our cluster:

```
export REDIS_PASSWORD=$(kubectl get secret --namespace default redis -o jsonpath="{.data.redis-password}" | base64 -d)
```

To check this, we can connect to the redis instance and run some queries:
First we port-forward the Redis Service port: 

```
kubectl port-forward svc/redis-master 6379:6379
```

Then we create a pod to be a redis client by running: 

```
kubectl run --namespace default redis-client --restart='Never'  --env REDIS_PASSWORD=$REDIS_PASSWORD  --image docker.io/bitnami/redis:6.2.5-debian-10-r63 --command -- sleep infinity
```

Then we attach to this pod

```
kubectl exec --tty -i redis-client --namespace default -- bash
```

Once we are in the pod, we can use the `redis-cli` to connect and send some queries:

```
redis-cli -h redis-master -a $REDIS_PASSWORD
```

You can run the following query to see which keys were stored: 

```
keys *
```

And then review the values by selecting a key with:

```
hgetall "generate-values||random-XX"
```

## Let's add tracing

```tracing.yaml
apiVersion: dapr.io/v1alpha1
kind: Configuration
metadata:
  name: tracing
  namespace: default
spec:
  tracing:
    samplingRate: "1"
    zipkin:
      endpointAddress: "http://zipkin.default.svc.cluster.local:9411/api/v2/spans"

```

Once again we need to choose an implementation, for this example we will install and use Zipkin, but you can change the underlaying implementation without the need to change any application code. 

```
kubectl create deployment zipkin --image openzipkin/zipkin
```

```
kubectl expose deployment zipkin --type ClusterIP --port 9411
```

## Let's enable and check metrics

```
kubectl create namespace dapr-monitoring
```

```
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm install dapr-prom prometheus-community/prometheus -n dapr-monitoring

```

