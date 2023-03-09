## Installation & pre-requisites

In this section we will: 
- Create a KinD Kubernetes Cluster
- Install Crossplane core components
- Install Crossplane Helm and Kubernetes providers
- Install Dapr

Let's create a KinD Cluster: 

```
kind create cluster
```

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
kubectl apply -f config/helm-provider-config.yaml
```

We also need to install the Crossplane Kubernetes Provider if we want to install custom resources. 

```
kubectl crossplane install provider crossplane/provider-kubernetes:main
```

Getting providers credentials
```
SA=$(kubectl -n crossplane-system get sa -o name | grep provider-kubernetes | sed -e 's|serviceaccount\/|crossplane-system:|g')
kubectl create clusterrolebinding provider-kubernetes-admin-binding --clusterrole cluster-admin --serviceaccount="${SA}"
```

And then configure it (this is only necessary if we are planning to install a kubernetes resource in the cluster where crossplane is installed):

```
kubectl apply -f config/kubernetes-provider-config.yaml
```

After a few seconds if you check the configured providers you should see both Helm and Kubernetes `INSTALLED` and `HEALTHY`: 

```
kubectl get providers.pkg.crossplane.io
NAME                             INSTALLED   HEALTHY   PACKAGE                               AGE
crossplane-provider-helm         True        True      crossplane/provider-helm:v0.10.0      49s
crossplane-provider-kubernetes   True        True      crossplane/provider-kubernetes:main   19s
```

Let's install Dapr next:

```
helm repo add dapr https://dapr.github.io/helm-charts/
helm repo update
helm upgrade --install dapr dapr/dapr \
--version=1.10.2 \
--namespace dapr-system \
--create-namespace \
--wait
```

Now we are ready to install our Databases Crossplane composition so we can provision all the components that our application needs to work.

## Working with Cloud Providers (example GCP)

If you want to extend the demo and create Cloud Resources in a Cloud Provider you need to install the corresponding Crossplane Provider, along side with the specific rights to create the resources that you are planning to provision using Crossplane resources.

We will use the Getting Started for GCP package that includes the GCP provider: 

```
kubectl crossplane install configuration registry.upbound.io/xp/getting-started-with-gcp:v1.10.2
```

Then we need to enable the services and create the roles needed for the provider to create specific cloud resources:

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

# enable cloud InMemoryStore API
SERVICE="redis.googleapis.com"
gcloud services enable $SERVICE --project $PROJECT_ID

# grant access to cloud InMemoryStore API
ROLE="roles/redis.admin"
gcloud projects add-iam-policy-binding --role="$ROLE" $PROJECT_ID --member "serviceAccount:$SA"

# create service account keyfile
gcloud iam service-accounts keys create creds.json --project $PROJECT_ID --iam-account $SA

```

Once we enabled the services and created the roles, we need to create a new Kubernetes Secret so the Crossplane Provider can access it: 

```
kubectl create secret generic gcp-creds -n crossplane-system --from-file=creds=./creds.json

```

Finally, we glue all the configurations together by creating a new `ProviderConfig` resource: 

```
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

Once you have the GCP provider configured you can deploy the compositions that create GKE resources, and use the `gke` level for matching selectors.

```
kubectl apply -f app-datahase-redis-gke.yaml
kubectl apply -f app-datahase-redis-gke-dapr.yaml
```

