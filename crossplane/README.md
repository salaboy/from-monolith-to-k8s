# Infrastructure for our Conference Application

In this step by step tutorial we will be using Crossplane to Provision the Redis and PostgreSQL instance for our application. 

By using Crossplane and Crossplane Compositions we are aiming to unify how these components are provisioned, hiding away where these components are for the end users (applicaiton teams).

Application teams should be able to request these resources using a declarative approach as with any other Kubernetes Resource. This enable teams to use Environment Pipelines to configure both the application services and the application infrastructure components needed by the application.

Make sure that you follow the [pre-requisites & installation](prerequisites.md) first.


## Databases on demand with Crossplane Compositions

First we will install a  Crossplane composition that uses the Crossplane Helm Provider to allow teams to request Databases on demand. 

```
kubectl apply -f databases/app-database-redis.yaml
kubectl apply -f databases/app-database-postgresql.yaml
kubectl apply -f databases/app-database-resource.yaml
```

The Crossplane Composition resource (`app-database-redis.yaml`) defines which cloud resources needs to be created and how they need to be configured together. The Crossplane Composite Resource Definition (XRD) defines a simplified interface that enable application development teams to easily request new databases by creating resources of this type.

# Let's provision a new Database

We can provision a new Database for our team to use by executing the following command: 

```
kubectl apply -f my-db.yaml
```

The `my-db-keyvalue.yaml` resource looks like this: 

```
apiVersion: salaboy.com/v1alpha1
kind: Database
metadata:
  name: my-db-keyvalue
spec:
  compositionSelector:
    matchLabels:
      provider: local
      type: dev
      kind: keyvalue
  parameters: 
    size: small
```

Notice that we are using the labels `provider: local`, `type: dev` and `kind: keyvalue`. This allows Crossplane to find the right composition based on the labels. In this case a local Redis instance created by the Helm Provider.

You can check the database status using:

```
> kubectl get dbs
NAME    SIZE    MOCKDATA   KIND       SYNCED   READY   COMPOSITION            AGE
my-db   small              keyvalue   True     True    db.local.salaboy.com   5s
```

You can check that a new Redis instance was created in the `my-db` namespace. 



