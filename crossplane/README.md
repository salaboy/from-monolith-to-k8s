# Crossplane for Infrastructure in GCP
This directory contains [Crossplane.io](http://crossplane.io) resources to expose and provision PostgreSQL and Redis in GCP by using Crossplane abstractions. 


## Installation

Install Crossplane Self-Hosted: https://crossplane.io/docs/v1.2/getting-started/install-configure.html

**Note**: Make sure that you also install the Crossplane Kubectl plugin. 

**Note**: The following enable Redis API, make sure you run the commands in the right order so the SA ends up with these Service and Role too. 

```
export SERVICE="redis.googleapis.com"
gcloud services enable $SERVICE --project $PROJECT_ID

export ROLE="roles/redis.admin"
gcloud projects add-iam-policy-binding --role="$ROLE" $PROJECT_ID --member "serviceAccount:$SA"

```
Then you can create the Service Account Key file, as instructed in the docs: 

```
# create service account keyfile
gcloud iam service-accounts keys create creds.json --project $PROJECT_ID --iam-account $SA
```

Once you have the credentials set, install the `fmtok8s` package into Crossplane. This Crossplane Package contains the `fmtok8s` abstractions for PostgreSQL and Redis. Example of these resources are located in the `/resources/` sub-directory. 


Inside the `pkg` directory run:

```
kubectl crossplane build configuration
```

And then you can push the generated package as an OCI image to any Container Registry of your choice, for example [DockerHub](hub.docker.com)

```
# Set this to the Docker Hub username or OCI registry you wish to use.
REG=salaboy
kubectl crossplane push configuration ${REG}/crossplane-fmtok8s-gcp:0.0.1
```

I've pushed mine here: [https://hub.docker.com/repository/docker/salaboy/crossplane-fmtok8s-gcp](https://hub.docker.com/repository/docker/salaboy/crossplane-fmtok8s-gcp)

Now you can distribute this package and install it in any Crossplane installation. 

To install this package in your Crossplane installation you need to run: 

```
kubectl crossplane install configuration salaboy/crossplane-fmtok8s-gcp:0.0.8
```

Once you have this running you can create a PostgreSQL database instance by running

```
kubectl apply -f resources/postgresql.yaml
```

And then modifying the C4P Service to use it by adding some environment variables:

```
        - name: SPRING_DATASOURCE_DRIVERCLASSNAME
          value: org.postgresql.Driver 
        - name: SPRING_DATASOURCE_PLATFORM
          value: postgres
        - name: SPRING_DATASOURCE_URL
          value: jdbc:postgresql://${DB_ENDPOINT}:${DB_PORT}/postgres
        - name: SPRING_DATASOURCE_USERNAME
          valueFrom:
            secretKeyRef:
              name: db-conn
              key: username
        - name: SPRING_DATASOURCE_PASSWORD
          valueFrom:
            secretKeyRef:
              name: db-conn
              key: password
        - name: DB_ENDPOINT
          valueFrom:
            secretKeyRef:
              name: db-conn
              key: endpoint
        - name: DB_PORT
          valueFrom:
            secretKeyRef:
              name: db-conn
              key: port
```

You can do the same for the Agenda Service with Redis: 

The variables needed for Redis are: 

```
- name: SPRING_REDIS_IN_MEMORY
  value: "false"
- name: SPRING_REDIS_HOST
  valueFrom: 
    secretKeRef:
      key: endpoint
      name: redis-conn
```
