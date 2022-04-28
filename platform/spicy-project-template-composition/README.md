# Crossplane Composition to create Worker Clusters

This directory contains the Crossplane.io resources to create new Kubernetes Clusters in GGP. This example can be extended to any other Cloud Provider. 

For this to work you need to have Crossplane installed in the Cluster and the GCP Provider. Follow the docs here: https://crossplane.io/docs/v1.7/getting-started/install-configure.html
When you install the GCP provider you need to add a couple of extra permissions for the GPC provider to create clusters, networks and nodepools:


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

and this: 

```
kubectl create secret generic gcp-creds -n crossplane-system --from-file=creds=./creds.json
```

## Creating Worker Clusters definitions

This direcotry contains 3 files that define a Crossplane Composition and a new CRD to create worker clusters. 

- crossplane.yaml
- composition.yaml
- definition.yaml

You can read more about how composition works here: https://crossplane.io/docs/v1.7/getting-started/create-configuration.html

I've packaged this composition as an OCI image, so you can just install it by running: 

```
kubectl crossplane install configuration salaboy/worker-cluster-gcp:0.1.0
```

You can modify and build the package yourself:
```
kubectl crossplane build configuration
kubectl crossplane push configuration <USER>/worker-cluster-gcp:0.1.0
```
