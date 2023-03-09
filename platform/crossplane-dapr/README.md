# Provisioning and Consuming Multi-Cloud Infrastructure with Crossplane and Dapr

This tutorial shows how we can combine Crossplane compositions and Dapr to abstract away the complexity of consuming application's infrastructure no matter where it is provisioned. 

For that we will install Crossplane and Dapr in our KinD cluster. 
Then we will create a Crossplane composition that uses the Crossplane Helm Provider to provision a Redis database for our application. 

Once we have all our application needs to work, we want to make the life of our developers easier, no matter which programming language they are using. For that we will use Dapr and the StateStore component to connect to the provisioned database. To simplify the whole journey we will add the Dapr Component configuration to the Crossplane composition, so whenever a developer request a new application context, they don't need to worry about where the Redis database is running or how to connect to it. 

You can read more about this tutorial in this blog post: <LINK>

Make sure that you follow the [pre-requisites installation first](pre-requisites.md)

## Databases on demand with Crossplane Compositions

First we will install a  Crossplane composition that uses the Crossplane Helm Provider to allow teams to request Databases on demand. 

```
kubectl apply -f app-database-redis.yaml
kubectl apply -f app-database-resource.yaml
```

The Crossplane Composition resource (`app-database-redis.yaml`) defines which cloud resources needs to be created and how they need to be configured together. The Crossplane Composite Resource Definition (XRD) defines a simplified interface that enable application development teams to easily request new databases by creating resources of this type.

Note: This resources be packaged as an OCI image using the crossplane `kubectl` plugin and distributed for other teams to consume.

# Let's provision a new Database

We can provision a new Database for our team to use by executing the following command: 

```
kubectl apply -f my-db.yaml
```

The `my-db.yaml` resource looks like this: 

```
apiVersion: salaboy.com/v1alpha1
kind: Database
metadata:
  name: my-db
spec:
  compositionSelector:
    matchLabels:
      provider: local
      type: dev
  parameters: 
    size: small
```

Notice that we are using the labels `provider: local` and `type: dev`. This allows Crossplane to find the right composition based on the labels. 

You can check the database status using:

```
> kubectl get dbs
NAME    SIZE    SYNCED   READY   COMPOSITION            AGE
my-db   small   True     True    db.local.salaboy.com   2m49s
```

You can check that a new Redis instance was created in the `my-db` namespace. 

Now that we can create Databases on demand, let's make them available to developers by using the Dapr Statestore component!

# Making infrastructure available to developers

We can extend our Crossplane Composition to use Dapr Components to enable developers to connect to Databases, Pub/Sub, Secret Stores, etc. without worrying about where these components are running or which libraries do they need to connect to them. 

With Dapr we can use the Statestore component to connect to our Redis instance and then by using the Dapr SDKs or plain HTTP/GRPC requests developers can store and read data without adding new dependencies. Even better, the Platform team can decide to change Redis by another Statestore supported implementation without impacting our application code. 

Let's install now a Crossplane Composition that also creates a Dapr Statestore component that is configured to connect to the Redis instance created by the Redis Helm Chart:

```
kubectl apply -f app-database-redis-dapr.yaml
```

Once this Crossplane Composition is installed we can create new databases using the new composition label `type: dapr-dev`:

```
apiVersion: salaboy.com/v1alpha1
kind: Database
metadata:
  name: my-db-dapr
spec:
  compositionSelector:
    matchLabels:
      provider: local
      type: dapr-dev
  parameters: 
    size: small
```


Apply this resource by running: 

```
kubectl apply -f my-dapr-db.yaml
```

Once again, you can check the state of these resources by running the following commands: 

```
> kubectl get dbs
NAME         SIZE    SYNCED   READY   COMPOSITION                 AGE
my-db        small   True     True    db.local.salaboy.com        21m
my-db-dapr   small   True     True    dapr.db.local.salaboy.com   29s
```

This composition created a new Dapr Statestore component inside the `my-db-dapr` database. You can check it by running: 

```
> kubectl describe components -n my-db-dapr my-db-dapr-statestore
Name:         my-db-dapr-statestore
Namespace:    my-db-dapr
API Version:  dapr.io/v1alpha1
Kind:         Component
Metadata:
  ...
Spec:
  Metadata:
    Name:  redisHost
    Secret Key Ref:
      Key:   url
      Name:  my-db-dapr-redis-dapr
    Name:    redisPassword
    Secret Key Ref:
      Key:   password
      Name:  my-db-dapr-redis-dapr
    Name:    keyPrefix
    Value:   name
  Type:      state.redis
  Version:   v1
Events:      <none>
```

As you can notice, both the `url` and the `password` to connect to the Redis instance are coming from a [secret that was also created by the Crossplane composition](https://github.com/salaboy/from-monolith-to-k8s/blob/main/platform/crossplane-dapr/app-database-redis-dapr.yaml#L54). 

# Connecting our applications to a Statestore component

In this section we will deploy two simple application. One in Java that writes data to the statestore and one in Go that reads the data from it. 

You can find the source code for these two simple applications in this repository [Dapr Example Apps](https://github.com/salaboy/dapr-example-apps). Notice that you don't need the source code to run this applications. 

The goal of this section is to highlight the polyglot appraoch of Dapr and show how without adding any Redis specific dependency to our applications, we can connect and use it. As this approach applies also to Pub/Sub, secret stores, workflows, actors, etc. 

Let's deploy our applications: 

`kubectl apply -f apps.yaml`

To interact with these two applications you will need to use `port-forward` to both services.

To access the Java Application using port 8080:
```
kubectl port-forward svc/java-app-service 8080:80 
```
To access the Go Application using port 8081:
```
kubectl port-forward svc/go-app-service 8081:80
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

