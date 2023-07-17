# Infrastructure for our Conference Application

In this step-by-step tutorial, we will be using Crossplane to Provision the Redis, PostgreSQL and Kafka instances for our application services to use. 

By using Crossplane and Crossplane Compositions, we aim to unify how these components are provisioned, hiding away where these components are for the end users (application teams).

Application teams should be able to request these resources using a declarative approach as with any other Kubernetes Resource. This enables teams to use Environment Pipelines to configure both the application services and the application infrastructure components needed by the application.

Make sure that you follow the [pre-requisites & installation](prerequisites.md) first.


## Databases on demand with Crossplane Compositions

First, we will install a Crossplane composition that uses the Crossplane Helm Provider to allow teams to request Databases on demand. 

```
kubectl apply -f databases/app-database-redis.yaml
kubectl apply -f databases/app-database-postgresql.yaml
kubectl apply -f databases/app-database-resource.yaml
```

The Crossplane Composition resource (`app-database-redis.yaml`) defines which cloud resources need to be created and how they need to be configured together. The Crossplane Composite Resource Definition (XRD) (`app-database-resource.yaml`) defines a simplified interface that enables application development teams to quickly request new databases by creating resources of this type.

## Let's provision a new Database

We can provision a new Database for our team to use by executing the following command: 

```
kubectl apply -f my-db-keyvalue.yaml
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

Notice that we are using the labels `provider: local`, `type: dev` and `kind: keyvalue`. This allows Crossplane to find the right composition based on the labels. In this case, a local Redis instance was created by the Helm Provider.

You can check the database status using:

```
> kubectl get dbs
NAME              SIZE    MOCKDATA   KIND       SYNCED   READY   COMPOSITION                     AGE
my-db-keyavalue   small   false      keyvalue   True     True    keyvalue.db.local.salaboy.com   97s
```

You can check that a new Redis instance was created in the `my-db-keyvalue` namespace. 

You can follow the same steps to provision a PostgreSQL database by running: 

```
kubectl apply -f my-db-sql.yaml
```

You should see now two `dbs`

```
> kubectl get dbs
NAME              SIZE    MOCKDATA   KIND       SYNCED   READY   COMPOSITION                     AGE
my-db-keyavalue   small   false      keyvalue   True     True    keyvalue.db.local.salaboy.com   2m
my-db-sql         small   false      sql        True     False   sql.db.local.salaboy.com        5s
```


You can now check that there are two Pods running, one for each database:

```
> kubectl get pods
NAME                             READY   STATUS    RESTARTS   AGE
my-db-keyavalue-redis-master-0   1/1     Running   0          3m40s
my-db-sql-postgresql-0           1/1     Running   0          104s
```

And there should be 4 Kubernetes Secrets (two for our two helm releases and two containing the credentials to connect to the newly created instances):

```
> kubectl get secret
NAME                                    TYPE                 DATA   AGE
my-db-keyavalue-redis                   Opaque               1      2m32s
my-db-sql-postgresql                    Opaque               1      36s
sh.helm.release.v1.my-db-keyavalue.v1   helm.sh/release.v1   1      2m32s
sh.helm.release.v1.my-db-sql.v1         helm.sh/release.v1   1      36s
```

## Let's deploy our Conference Application

Ok, now that we have our two databases running, we need to make sure that our application services connect to these instances. The first thing that we need to do is to disable the Agenda and Call For Proposal Services helm dependencies so that when the charts get installed don't install new databases. 

For that, we will use the `app-values.yaml` file containing the configurations for the services to connect to our newly created databases:

```
helm repo add fmtok8s https://salaboy.github.io/helm/
helm repo update
helm install conference fmtok8s/fmtok8s-conference-chart -f app-values.yaml
```

The `app-values.yaml` content looks like this: 
```
fmtok8s-agenda-service: 
  redis:
    enabled: false
  env: 
    - name: SPRING_REDIS_HOST
      value: my-db-keyavalue-redis-master
    - name: SPRING_REDIS_PORT
      value: "6379" 
    - name: SPRING_REDIS_PASSWORD
      valueFrom:
        secretKeyRef:
          name: my-db-keyavalue-redis
          key: redis-password
    
fmtok8s-c4p-service: 
  postgresql:
    enabled: false
  env: 
  - name: DB_ENDPOINT
    value: my-db-sql-postgresql
  - name: DB_PORT
    value: "5432"
  - name: SPRING_R2DBC_PASSWORD
    valueFrom:
      secretKeyRef:
        name: my-db-sql-postgresql
        key: postgres-password
```

As you can see, we are just setting Environment Variables and referencing the secrets for the database passwords. 

## Sum up

In this tutorial, we have managed to separate the provisioning for the application infrastructure from the application deployment. This enables different teams to request resources on-demand (using Crossplane compositions) and application services that can evolve independently. 

Using Helm Chart dependencies for development purposes and quickly getting a fully functional instance of the application up and running is great. For more sensitive environments you might want to follow an approach like the one shown here, where you have multiple ways to connect your application with the components required by each service. 

