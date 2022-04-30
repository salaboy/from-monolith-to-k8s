# Creating your own self-service platform

With a fast-paced CNCF landspace, it is a full time job to understand, pick and glue projects together to enable your teams with a self-service platform to suit their needs. 

In this document we will be using Kratix, Crossplane, Knative Serving and Knative Functions to demonstrate how a Platform Team can build and curate a set of tools that will enable teams to request using a declarative way. This is just an example of how you can achieve this, and other tools can be used to implement the same behaviours, but some key points that we have tried to cover here are: 
- Self-Service Platform covering two main Personas: Platform Engineers (Platform Team) , Developer Teams (App Team) 
- Developer experience is key to improve productivity, reducing the cognitive load for teams is key to improve efficiency
- Platform Teams can collaborate with the teams to create the right platform for them, using exensible mechanism that can be adapted for more complex needs when needed
- You can achieve all of this by using Open Source projects, but you will need to provide your domain-specific glue


In general, the Platform Team will be in charge of creating the "Platform", developer teams should interact with the "Platform" probably using a portal to create their requests and obtain their environment's credentials. 

![Platform Teams and Developer Teams]


## Use Case
 
(TBD)


## Creating a Self-Service Platform for K8s

Because we are using Kubernetes as the target platform, we will be using [Kratix](https://github.com/syntasso/kratix) to provide a consistent API using Kubernetes resources to enable developer teams to create their requests to the platform. 
This means that our Platform will be materialized as a Kubernetes Cluster with Kratix installed in it. Kratix, besides exposing a consistent and high-level API, will be in charge of orchestrating different tools inside this cluster to provision team environments. 

I am running this demo in GCP GKE, but you can run this in any Cloud Provider of your choice. 
In this folling sections, we will install all we need in the Platform Cluster, hence you can go ahead and create a Kubernetes Cluster, connect to it and then let's install some components: 
- Installing Kratix
- Installing Crossplane and the GCP provider
- Understanding what we have installed and why

### Installing Kratix into the Platform cluster

First it is recommended for you to clone Kratix repository so then we can modify and apply some Kubernetes resources:

```
git clone https://github.com/syntasso/kratix
```

**Note**: I am using the `dev` branch 

Before installing Kratix, I need to modify the following file `hack/platform/minio-install.yaml` to make sure that the Minio Service is exposed outside the cluster, hence I'm changing the `ServiceType` to `LoadBalancer`

Now let's install Kratix: 

```
kubectl apply -f distribution/kratix.yaml
kubectl apply -f hack/platform/minio-install.yaml
```

We are using Minio inside the Cluster to store resources that will be applied to environment clusters later on. 

**Note**: Check the Minio Service External IP with `kubectl get svc/minio -n kratix-platform-system` .  

![Kratix Diagram]()

#### Accessing Minio files from local environment
If you want to access from your local environment then use `Port Forward`

```
k port-forward svc/minio -n kratix-platform-system 9000:80
```


Then configure local `minio-mc` client to connect: 
```
vi /Users/msalatino/.mc/config.json

```

Configure local keys to `minioamdin`/`minioamdin`

Check Minio, this creates the CRDs inside Minio for the Worker Cluster (using flux) to pick up and apply: 

```
mc ls local/kratix-crds
```


###  Installing Crossplane into the Platform Cluster

We will be installing Crossplane into the platform cluster so the platform can provision new environments for the teams. We don't want to expose Crossplane Resources to our development teams, hence Crossplane resources will be applied by Kratix. 

In this section we will install Crossplane in the Platform Cluster and we need to make sure that Crossplane have enough rights to provision new Kubernetes Clusters for our teams. 

Check the Crossplane [installation instructions for GCP here](https://crossplane.io/docs/v1.7/getting-started/install-configure.html)
You can also install the `kubectl` plugin if you want to play around with the Crossplane composition that I've created for this example. 

When installing Crossplane you need to make sure that you give Crossplane permissions to create new GKE clusters and containers. You will need to enable the following services and create the following roles before creating the GCP credentials: 

```
SERVICE="container.googleapis.com"
gcloud services enable $SERVICE --project $PROJECT_ID
```

```
ROLE="roles/container.admin"
gcloud projects add-iam-policy-binding --role="$ROLE" $PROJECT_ID --member "serviceAccount:$SA"
```

```
SERVICE="compute.googleapis.com"
gcloud services enable $SERVICE --project $PROJECT_ID
```

```
ROLE="roles/compute.admin"
gcloud projects add-iam-policy-binding --role="$ROLE" $PROJECT_ID --member "serviceAccount:$SA"
```

Make sure that you do this before creating the service account secret, meaning before you do this: 
```

# create service account keyfile
gcloud iam service-accounts keys create creds.json --project $PROJECT_ID --iam-account $SA

```

and then you can create the secret with the credentials in it: 

```
kubectl create secret generic gcp-creds -n crossplane-system --from-file=creds=./creds.json
```

By doing this, now Crossplane will be able to create GKE clusters on your behalf. 


## 

## Reference
- [Kratix Getting Started Guide](https://github.com/syntasso/kratix/blob/main/docs/quick-start.md)
- [Crossplane and GCP Provider](https://crossplane.io/docs/v1.7/getting-started/install-configure.html)
- 
