# Polyglot CloudEvents consumer/producer Tutorial

In this short tutorial we are going to create two applications that produce and consume CloudEvents. These applications uses different technology stacks: 
- Application A ([`fmtok8s-java-cloudevents`](https://github.com/salaboy/fmtok8s-java-cloudevents)) uses Java and Spring Boot and it adds the CloudEvents Java SDK to write and read CloudEvents. 
- Application B ([`fmtok8s-go-cloudevents`](https://github.com/salaboy/fmtok8s-go-cloudevents)) uses Go and adds the CloudEvents Go SDK to read and write CloudEvents. 

![CloudEvents Examples](cloudevents-fmtok8s.png)

If you want to build and run the applications you will need to hava Java, Maven and Go installed. Alternatively you can run the available Docker containers, in which case you only need Docker to run this examples. 

To consume CloudEvents via HTTP both applications expose a REST endpoint where the events can be received.

## Running the applications in your local environment
You can find the source code of each application in the following repositories:

- Application A ([`fmtok8s-java-cloudevents`](https://github.com/salaboy/fmtok8s-java-cloudevents)) uses Java and Spring Boot and it adds the CloudEvents Java SDK to write and read CloudEvents. 
- Application B ([`fmtok8s-go-cloudevents`](https://github.com/salaboy/fmtok8s-go-cloudevents))

You can build the application from source using Java and Go tooling (as explained in each repository), but if you want to get things up and running fast I recommend just using the available docker containers. 

Because containers will be talking to each other, you need to make sure that they are in the same docker network, for that reason, we will create a custom network (`cloudevents-net`) and make sure that both containers are attached to it: 

```
docker network create --driver bridge cloudevents-net
```

If you have docker installed in your environment run the following commands in two separate terminals: 

Application A (Java):
```
docker run --name application-a --network cloudevents-net -e SINK=http://application-b:8081 -p 8080:8080 salaboy/fmtok8s-java-cloudevents
```

Application B (Go):
```
docker run --name application-b --network cloudevents-net -e SINK=http://application-a:8080 -p 8081:8081 salaboy/fmtok8s-go-cloudevents.go
```

![CloudEvents Examples Docker](cloudevents-fmtok8s-docker.png)

Notice that it is pretty mucht the same command for both applications, just changing the SINK for the produced event and the application name. It is important to notice here that we are setting the container's name so they can call each other using a fixed name instead of an IP address. 

From a third terminal you can ask Application A to produce a CloudEvent which will be sent to Application B (as the SINK is configured to point to Application B port on localhost 8081)

You can use `curl` to send a POST request to the `/produce` endpoint of Application A (localhost:8080)
```
curl -X POST http://localhost:8080/produce
```

Check the logs from both containers by looking into the tabs where you run the previous commands. 

This POST request will hit Application A which will produce a new CloudEvent, this CloudEvent will be sent to Application B `/` endpoint which is ready to consume a CloudEvent and print it's contants after parsing the body into a Go struct. 

Feel free to test the other way around, by calling the `/produce` endpoint in Application B. 

```
curl -X POST http://localhost:8081/produce
```

You can also curl with CloudEvents to each application, for `application-a`
```
curl -X POST http://localhost:8080/ -H "Content-Type: application/json" -H "ce-type: app-b.MyCloudEvent"  -H "ce-id: 123"  -H "ce-specversion: 1.0" -H "ce-source: curl-command" -d '{"myData" : "hello from curl", "myCounter" : 1 }'
```
Same for `application-b` just different port: 

```
curl -X POST http://localhost:8081/ -H "Content-Type: application/json" -H "ce-type: app-a.MyCloudEvent"  -H "ce-id: 123"  -H "ce-specversion: 1.0" -H "ce-source: curl-command" -d '{"myData" : "hello from curl", "myCounter" : 1 }'
```

# On Kubernetes

To run the same services inside Kubernetes you just need to have the right Kubernetes resources. You can find two YAML files inside the [`kubernetes/`](https://github.com/salaboy/from-monolith-to-k8s/tree/master/cloudevents/kubernetes) directory. These YAML files contains a Kubernetes Deployment and a Kubernetes Service definition for each service. 
By deploying these two services (`application-a-service` and `application-b-service`) on Kubernetes we are not changing the topology or the fact that one services needs to know the other service name in order to send a CloudEvent. 

![CloudEvents Examples Kubernetes](cloudevents-fmtok8s-kubernetes.png)

You will notices inside the Kubernetes Deployment of both applications that we are defining the SINK variable as we were doing with Docker. When deploying inside Kubernetes and using Kubernetes Services, we can use the Service name to interact with our containerized applications. Notice that in Kubernetes, we don't need to create any new network (as required with Docker) to be able to use the Service name discovery mechanism. By using the service name, we rely on Kubernetes to route the traffic to the right container. 

Running the following command you can get both applications up and running in your Kubernetes Cluster:
```
kubectl apply -f kubernetes/
```

Now if you want to interact with the services that are now running inside Kubernetes you can use `port-forward` for `application-a`: 

```
kubectl port-forward svc/application-a-service 8080:80
```

and for `application-b`:
```
kubectl port-forward svc/application-b-service 8081:80
```

Now all the traffic that you send to `localhost:8080` or `localhost:8081` will go to `application-a` and `application-b` respectively. 

You can try the produce endpoint on `application-a`, as we did with Docker:
```
curl -X POST http://localhost:8080/produce
```
or to `application-b`:

```
curl -X POST http://localhost:8081/produce
```

You can inspect the logs of both applications by using `kubectl logs -f <POD_NAME>`.

Same you can send CloudEvents directly to each application, for example to `application-a`: 
```
curl -X POST http://localhost:8080/ -H "Content-Type: application/json" -H "ce-type: app-b.MyCloudEvent"  -H "ce-id: 123"  -H "ce-specversion: 1.0" -H "ce-source: curl-command" -d '{"myData" : "hello from curl", "myCounter" : 1 }'
```

and to `application-b` is the same, just different port: 
```
curl -X POST http://localhost:8081/ -H "Content-Type: application/json" -H "ce-type: app-a.MyCloudEvent"  -H "ce-id: 123"  -H "ce-specversion: 1.0" -H "ce-source: curl-command" -d '{"myData" : "hello from curl", "myCounter" : 1 }'
```

# With Knative Eventing
So far, applications are sending Events to each other, but if we are building Event-Driven applications we might want to decouple producers from consumers. 
To achieve this more decoupled architecture we will use Knative Eventing. 
If you have a Kubernetes Cluster you can install [Knative Eventing by following the getting started guide in the official site](https://knative.dev/docs/install/eventing/install-eventing-with-yaml/).

Make sure that you install the **In Memory Standalone channel** and the **MT-Channel-Based broker**, which both are listed optional in the installation guide.

We will be deploying the same applications that we had deployed in the previous steps but they will not be sending events to each other directly. Instead, each application will know only about an Event Broker.

![Knative Eventing into the mix](cloudevents-fmtok8s-with-knative-eventing.png)

Once we have Knative Eventing installed and the Channel and Broker implementation, we need to create a Broker instance for our applications to use. To create a broker run: 
```
kubectl create -f - <<EOF
apiVersion: eventing.knative.dev/v1
kind: Broker
metadata:
 name: default
 namespace: default
EOF
``` 

Check that the broker is in Ready state and that it provides an URL: 

````
salaboy> kubectl get broker
NAME      URL                                                                        AGE   READY   REASON
default   http://broker-ingress.knative-eventing.svc.cluster.local/default/default   2s    True   
```


The only change that we need to perform in our application's configuration is the `SINK` environment variable which now need to point to the Knative Broker that we just created. Now your applications can send events to the broker URL (http://broker-ingress.knative-eventing.svc.cluster.local/default/default).

If you have the applications already running in the cluster change the `SINK` variable value to `http://broker-ingress.knative-eventing.svc.cluster.local/default/default`. You can do that by editing the application's deployments: 

```
kubectl edit deploy application-a
```

Do the same for `application-b`.

Notice that now, both of the applications will be sending events to the Broker, but the broker will not forward these events, because we haven't created any subscription to them. 

Let's start by creating a Knative trigger (subscription) for `application-b` to receive CloudEvents sent to the Broker: 

```
kubectl create -f - <<EOF
apiVersion: eventing.knative.dev/v1
kind: Trigger
metadata:
  name: app-b-trigger
  namespace: default
spec:
  broker: default
  subscriber:
    ref:
      apiVersion: v1
      kind: Service
      name: application-b-service
EOF
```

Check that the trigger was created and it is ready:

```
salaboy> kubectl get trigger
NAME            BROKER    SUBSCRIBER_URI                 AGE   READY   REASON
app-b-trigger   default   http://application-b-service   3s    True
```

Now if you produce an event from `application-a` the event will be sent to the Broker and the recently created Trigger will forward the event to `application-b`.

Make sure that you have access to `application-a` by using `port-foward`
```
kubectl port-forward svc/application-a-service 8080:80
```

And then produce a cloud event by using `curl`:

```
curl -X POST http://localhost:8080/produce
```

If you check the logs (`kubectl logs -f <POD_NAME>`) of `application-b` you should see the CloudEvent arriving via the Broker. 

Now if you add a trigger for `application-a`, as follow: 

```
kubectl create -f - <<EOF
apiVersion: eventing.knative.dev/v1
kind: Trigger
metadata:
  name: app-b-trigger
  namespace: default
spec:
  broker: default
  subscriber:
    ref:
      apiVersion: v1
      kind: Service
      name: application-b-service
EOF
```

You will notice that no matter which application produces an event, both applications will get it. Also notice that you need to use the full service name, including the namespace and `svc.cluster.local` for the event to arrive correctly. 

For this example, `apiVersion: v1` and `kind: Service` makes reference to Kubernetes Services. But if you are using Knative Serving Services you will need to use `apiVersion: serving.knative.dev/v1` and `kind: Service`. 


## Knative Eventing Extras

This section covers a couple of details that you need to be aware when using Knative Eventing
- `REF` vs `URI` nd oher properties for subscribers in Knative Eventing Triggers
- Knative Eventing Trigger filters
- Monitoring Events using Sockeye

### `REF` vs `URI` and oher properties for subscribers in Knative Eventing Triggers

When defining triggers, in the previous example, to point to a subscriber we have used `ref` which allows us to point to an `Addressable` Kubernetes resource. This allows Knative Eventing to validate that the resource is there to send events to. The trigger will fail to be ready if the resource is not present. 

But there are situations where you don't have an addressable resource or you want to avoid Knative Eventing checking on the subscriber. For this scenarios you can use `uri`. In our previous example it will look as follows:  

```
kubectl create -f - <<EOF
apiVersion: eventing.knative.dev/v1
kind: Trigger
metadata:
  name: app-a-trigger
  namespace: default
spec:
  broker: default
  subscriber:
    uri: http://application-a-service.default.svc.cluster.local
EOF
```

Another important detail to know is that you can use `uri` in conjunction with `ref` to point to path different from `/`. For example, if you are listening for CloudEvents under `/events` you can use `uri` to add `/events` on top of a `ref` resource, as shown below:

```
kubectl create -f - <<EOF
apiVersion: eventing.knative.dev/v1
kind: Trigger
metadata:
  name: app-b-trigger
  namespace: default
spec:
  broker: default
  subscriber:
    ref:
      apiVersion: v1
      kind: Service
      name: application-a-service  
    uri: /events
EOF
```

### Knative Eventing Trigger Filters
As covered in [the official documentation](https://knative.dev/docs/eventing/broker/triggers/#trigger-filtering) you can apply filters to CloudEvent attributes, including extensions. For now, only exact matches on strings are supported, the community is currently working on adding more advanced techniques for filtering events. 

It is quite common for your applications to be interested in certain types of events and you can easily achieve this by adding the filter attribute to your triggers, as shown below: 

```
kubectl create -f - <<EOF
apiVersion: eventing.knative.dev/v1
kind: Trigger
metadata:
  name: app-b-trigger
  namespace: default
spec:
  broker: default
  filter:
    attributes:
      type: MyCloudEvent
  subscriber:
    ref:
      apiVersion: v1
      kind: Service
      name: application-a-service  
EOF
```
### Monitoring Events using Sockeye
Finally, it is always good to tap into all events to monitor what is going on with your producers and consumers. This can be easily achieved by creating a trigger without any filters and redirecting them to a service or application that allows you to see every event that is sent to the Broker.

One of these applications that you can use is called Sockeye and it has been used a lot in the Knative Eventing community. Don't expect anything complicated it is just a web applicaiton that will print every CloudEvent that receives. 

You can install Sockeye if you have Knative Serving installed with: 

```
kubectl apply -f https://github.com/n3wscott/sockeye/releases/download/v0.7.0/release.yaml
```

And then create a trigger to forward all the events from the Broker to Sockeye:
```
kubectl create -f - <<EOF
apiVersion: eventing.knative.dev/v1
kind: Trigger
metadata:
  name: wildcard-trigger
  namespace: default
spec:
  broker: default
  subscriber:
    uri: http://sockeye.default.svc.cluster.local
EOF
```
