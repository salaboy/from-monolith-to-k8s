# Installing Application Infrastructure with Helm

This section covers how to install Application Infrastructure with Helm and how to connect these components (databases) with our existing Services. 

If you have a Kubernetes Cluster already with the [application services up and running](https://github.com/salaboy/from-monolith-to-k8s/tree/master/helm) you can proceed to install both PostgreSQL and Redis with Helm.


```
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
```

Install PostgreSQL Chart

```
helm install postgresql bitnami/postgresql
```

This should show an output similar to the following: 

```
NAME: postgresql
LAST DEPLOYED: Fri Jul 23 11:55:36 2021
NAMESPACE: default
STATUS: deployed
REVISION: 1
TEST SUITE: None
NOTES:
** Please be patient while the chart is being deployed **

PostgreSQL can be accessed via port 5432 on the following DNS names from within your cluster:

    **postgresql.default.svc.cluster.local - Read/Write connection**

To get the password for "postgres" run:

    **export POSTGRES_PASSWORD=$(kubectl get secret --namespace default postgresql -o jsonpath="{.data.postgresql-password}" | base64 --decode)**

To connect to your database run the following command:

    kubectl run postgresql-client --rm --tty -i --restart='Never' --namespace default --image docker.io/bitnami/postgresql:11.12.0-debian-10-r44 --env="PGPASSWORD=$POSTGRES_PASSWORD" --command -- psql --host postgresql -U postgres -d postgres -p 5432



To connect to your database from outside the cluster execute the following commands:

    kubectl port-forward --namespace default svc/postgresql 5432:5432 &
    PGPASSWORD="$POSTGRES_PASSWORD" psql --host 127.0.0.1 -U postgres -d postgres -p 5432
```

Highlighted in bold are the service DNS name that you will need to use to connect to the DB instance and how to obtain the password which is stored inside a Kubernetes Secret. You don't need to access the password in this way, you just need to know that a Kubernetes Secret has the password stored for your Pods to use when trying to connect to it. 

Then install Redis: 

```
helm install redis bitnami/redis
```

You should see the following output: 

```
NAME: redis
LAST DEPLOYED: Fri Jul 23 11:58:18 2021
NAMESPACE: default
STATUS: deployed
REVISION: 1
TEST SUITE: None
NOTES:
** Please be patient while the chart is being deployed **

Redis(TM) can be accessed on the following DNS names from within your cluster:

    **redis-master.default.svc.cluster.local for read/write operations (port 6379)
    redis-replicas.default.svc.cluster.local for read-only operations (port 6379)**



To get your password run:

    **export REDIS_PASSWORD=$(kubectl get secret --namespace default redis -o jsonpath="{.data.redis-password}" | base64 --decode)**

To connect to your Redis(TM) server:

1. Run a Redis(TM) pod that you can use as a client:

   kubectl run --namespace default redis-client --restart='Never'  --env REDIS_PASSWORD=$REDIS_PASSWORD  --image docker.io/bitnami/redis:6.2.5-debian-10-r0 --command -- sleep infinity

   Use the following command to attach to the pod:

   kubectl exec --tty -i redis-client \
   --namespace default -- bash

2. Connect using the Redis(TM) CLI:
   redis-cli -h redis-master -a $REDIS_PASSWORD
   redis-cli -h redis-replicas -a $REDIS_PASSWORD

To connect to your database from outside the cluster execute the following commands:

    kubectl port-forward --namespace default svc/redis-master 6379:6379 &
    redis-cli -h 127.0.0.1 -p 6379 -a $REDIS_PASSWORD
```

Same here, highlighted in bold the Service URL and the password secret. 

Now if you check the pods that are running besides your application you should see the following:

```
> kubectl get pods
NAME                                       READY   STATUS      RESTARTS   AGE
app-fmtok8s-agenda-rest-68bd9c8bcb-dxbdt   1/1     Running     0          12m
app-fmtok8s-api-gateway-58d49588b4-vw7pc   1/1     Running     0          12m
app-fmtok8s-c4p-rest-7cb8bc4485-scj7r      1/1     Running     0          12m
app-fmtok8s-email-rest-8f954fbbd-p8hld     1/1     Running     0          12m
postgresql-postgresql-0                    1/1     Running     0          6m30s
redis-master-0                             1/1     Running     0          3m47s
redis-replicas-0                           1/1     Running     0          3m47s
redis-replicas-1                           1/1     Running     0          3m13s
redis-replicas-2                           1/1     Running     0          2m36s
```

Both PostgreSQL and a Redis Cluster are running and ready to be used by our application. 

Now it is time to connect our services to these databases, let's start with the Call for Proposal Sevice (C4P). 
Edit the deployment called `app-fmtok8s-c4p-rest` and then look for the section `spec.template.spec.containers[0].env` and add the following variables:


```
        - name: SPRING_DATASOURCE_DRIVERCLASSNAME
          value: org.postgresql.Driver
        - name: SPRING_DATASOURCE_PLATFORM
          value: org.hibernate.dialect.PostgreSQLDialect
        - name: SPRING_DATASOURCE_URL
          value: jdbc:postgresql://${DB_ENDPOINT}:${DB_PORT}/postgres
        - name: SPRING_DATASOURCE_USERNAME
          value: postgres
        - name: SPRING_DATASOURCE_PASSWORD
          valueFrom:
            secretKeyRef:
              key: postgresql-password
              name: postgresql
        - name: DB_ENDPOINT
          value: postgresql
        - name: DB_PORT
          value: "5432"

```

Now you can scale up the replicas of your Call for Proposals (C4P) services and all the instances will connect to the same database. You have made your Call for Proposal Service stateless. 

You can do exactly the same with the Agenda Service and Redis. Let's edit the deployment and add the following new variables in the `spec.template.spec.containers[0].env` section: 

```
        - name: SPRING_REDIS_IN_MEMORY
          value: "false"
        - name: SPRING_REDIS_HOST
          value: redis-master
        - name: SPRING_REDIS_PASSWORD
          valueFrom:
            secretKeyRef:
              key: redis-password
              name: redis

```

Notice that for ready we need to explicitely disable the in-memory db that is running inside our service. 




This approach is recommended for experimenting, development and maybe testing, but you should check [Crossplane](../crossplane/README.md) for production usage in Cloud Providers. This is mostly you will need to maintain these components in the long run, including upgrading versions, backing up data, etc.

 
