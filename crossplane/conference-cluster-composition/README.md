# Creating a Cluster with Conference Platform inside

**THIS IS STILL WORK IN PROGRESS, due to an existing Bug with GKE**

In this case we are not ok with just creating a namepsace for our different Customers, we need a full blown Kubernetes Cluster for each application instance. 

This requires to setup the provider-gcp to also have this service and role:

```
SERVICE="container.googleapis.com"
gcloud services enable $SERVICE --project $PROJECT_ID
```

```
ROLE="roles/container.admin"
gcloud projects add-iam-policy-binding --role="$ROLE" $PROJECT_ID --member "serviceAccount:$SA"

```

Then inside this cluster we can install our conference platform in complete isolation with other customers. 
 