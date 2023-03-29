# Connecting our applications to a Statestore component

In this section we will deploy two simple application. One in Java that writes data to the statestore and one in Go that reads the data from it. 

The goal of this section is to highlight the polyglot appraoch of Dapr and show how without adding any Redis specific dependency to our applications, we can connect and use it. As this approach applies also to Pub/Sub, secret stores, workflows, actors, etc. 

Let's deploy our applications: 

`kubectl apply -f apps.yaml -n <namespace where the statestore component is>`

To interact with these two applications you will need to use `port-forward` to both services.

To access the Java Application using port 8080:
```
kubectl port-forward svc/java-app-service 8080:80 -n <namespace>
```
To access the Go Application using port 8081:
```
kubectl port-forward svc/go-app-service 8081:80 -n <namespace>
```

Now you can send request to both services. If you send first a request to the Go service that reads data from Redis you will see an empty array back. 

```
> curl localhost:8081          
{"Values":null} 
```

Let's now send a request to the Java application to store a value: 

```
> curl -X POST "localhost:8080/?value=42"
{"values":["42"]}
```

You can try adding another value. Notice that the request returns the currently stored values too.

Then the Go Application can return the stored value: 

```
> curl localhost:8081
{"Values":["42"]}
```

If you want to delete all values you can send a `DELETE` request to the Java Application: 

```
> curl -X DELETE localhost:8080
```
