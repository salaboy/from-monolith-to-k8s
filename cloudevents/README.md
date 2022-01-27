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

# On Kubernetes

To run the same services inside Kubernetes you just need to have the right Kubernetes resources to run these containers. You can find two YAML files inside the `kubernetes` directory. These YAML files contains a Kubernetes Deployment and a Kubernetes Service definition for each service. 
By deploying these two services (application-a and application-b) on Kubernetes we are not changing the topology or the fact that one services needs to know the other service name in order to send a CloudEvent. 

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

# With Knative Eventing
(TBD)
