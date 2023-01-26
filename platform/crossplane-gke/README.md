# Crossplane Environment Composition for platform building


This tutorial focus on GCP, but new resources can be created for AKS and EKS. 

This Crossplane Composite resource creates the following resources:
- Kubernetes Cluster (GCP Cluster)
- NodePool
- Helm Provider Config
- Helm Release inside the created cluster

## Installation

In a GCP Kubernetes Cluster install the following components.


Let's install [Crossplane](https://crossplane.io) into its own namespace using Helm: 

```

helm repo add crossplane-stable https://charts.crossplane.io/stable
helm repo update

helm install crossplane --namespace crossplane-system --create-namespace crossplane-stable/crossplane --wait
```

Install the `kubectl crossplane` plugin: 

```
curl -sL https://raw.githubusercontent.com/crossplane/crossplane/master/install.sh | sh
sudo mv kubectl-crossplane /usr/local/bin
```

Then install the Crossplane Helm provider: 
```
kubectl crossplane install provider crossplane/provider-helm:v0.10.0
```

We need to get the correct ServiceAccount to create a new ClusterRoleBinding so the Helm Provider can install Charts on our behalf. 

```
SA=$(kubectl -n crossplane-system get sa -o name | grep provider-helm | sed -e 's|serviceaccount\/|crossplane-system:|g')
kubectl create clusterrolebinding provider-helm-admin-binding --clusterrole cluster-admin --serviceaccount="${SA}"
```

```
kubectl apply -f crossplane/config/helm-provider-config.yaml
```

## Install & Configure GCP provider

I've used the GCP getting started, but GCP provider can be installed in the same way. 


```
kubectl crossplane install configuration registry.upbound.io/xp/getting-started-with-gcp:v1.10.2
```

Configure GCP to enable the provider to create resources on our behalf:

```
# replace this with your own gcp project id and the name of the service account
# that will be created.
PROJECT_ID=my-project
NEW_SA_NAME=test-service-account-name

# create service account
SA="${NEW_SA_NAME}@${PROJECT_ID}.iam.gserviceaccount.com"
gcloud iam service-accounts create $NEW_SA_NAME --project $PROJECT_ID

# enable cloud SQL API
SERVICE="sqladmin.googleapis.com"
gcloud services enable $SERVICE --project $PROJECT_ID

# grant access to cloud SQL API
ROLE="roles/cloudsql.admin"
gcloud projects add-iam-policy-binding --role="$ROLE" $PROJECT_ID --member "serviceAccount:$SA"

# enable cloud container API
SERVICE="container.googleapis.com"
gcloud services enable $SERVICE --project $PROJECT_ID

# grant access to cloud container API
ROLE="roles/container.admin"
gcloud projects add-iam-policy-binding --role="$ROLE" $PROJECT_ID --member "serviceAccount:$SA"

# enable cloud compute API
SERVICE="compute.googleapis.com"
gcloud services enable $SERVICE --project $PROJECT_ID

# grant access to cloud compute API
ROLE="roles/compute.admin"
gcloud projects add-iam-policy-binding --role="$ROLE" $PROJECT_ID --member "serviceAccount:$SA"

# create service account keyfile
gcloud iam service-accounts keys create creds.json --project $PROJECT_ID --iam-account $SA

```

`PROJECT_ID` is your GCP project

```
kubectl create secret generic gcp-creds -n crossplane-system --from-file=creds=./creds.json

```

Finally

```
# replace this with your own gcp project id
PROJECT_ID=my-project
echo "apiVersion: gcp.crossplane.io/v1beta1
kind: ProviderConfig
metadata:
  name: default
spec:
  projectID: ${PROJECT_ID}
  credentials:
    source: Secret
    secretRef:
      namespace: crossplane-system
      name: gcp-creds
      key: creds" | kubectl apply -f -

```

## Install Environment Composite Resource (XRD)

```
kubectl apply -f environment-composition.yaml
kubectl apply -f environment-resource-definition.yaml
```

This can be packaged as an OCI image using the crossplane `kubectl` plugin.

# Let's provision new Environments

```
kubectl apply -f team-a-env.yaml
```


You can check the environment status using:

```
kubectl get env
```

Once the cluster is created you can connect to it using `gcloud connect...` from the GCP dashboard and check that the `dapr-system` is installed and working. 

```
gcloud connect ..
kubectl get pods -n dapr-system
```

