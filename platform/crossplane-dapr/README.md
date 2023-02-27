# Provisioning and Consuming Multi-Cloud Infrastructure with Crossplane and Dapr

This tutorial shows how we can combine Crossplane compositions and Dapr to abstract away the complexity of consuming application's infrastructure no matter where it is provisioned. 

For that we will install Crossplane and Dapr in our KinD cluster. 
Then we will create a Crossplane composition that uses the Crossplane Helm Provider to provision a Redis database for our application. 

Once we have all our application needs to work, we want to make the life of our developers easier, no matter which programming language they are using. For that we will use Dapr and the StateStore component to connect to the provisioned database. To simplify the whole journey we will add the Dapr Component configuration to the Crossplane composition, so whenever a developer request a new application context, they don't need to worry about where the Redis database is running or how to connect to it. 



## Installation

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

Now we are ready to install our Development Environment Crossplane composition so we can provision all the components that our application needs to work.

## Install Environment Composite Resource (XRD)

```
kubectl apply -f env-composition-gke.yaml
kubectl apply -f env-resource.yaml
```

This can be packaged as an OCI image using the crossplane `kubectl` plugin.

# Let's provision new Environments

```
kubectl apply -f env.yaml
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


# VCluster 

Without installing anything else, you can provision a new environment using VCluster instead of GKE. 

For this you only need the VCluster composition: 

```
kubectl apply -f env-composition-vcluster.yaml
```

Now that the composition is ready, you can provision a new environment for `team-b`
```
kubectl apply -f team-b-env-vcluster.yaml
```

When the environment is ready you can connect using the `vcluster` CLI: 

```
vcluster connect team-b-env --server https://localhost:8443 -- zsh
```

Now you are in your VCluster provisioned by Crossplane! Check that there are some pods in the default namespace for the helm chart that the composition installed: 

```
kubectl get pods 
```


Notice that `team-b-env-vcluster.yaml` and `team-a-env-gke.yaml` are exactly the same besides the `matchLabels` that allows Crossplane to pick the right composition depending which provider do we want to use.


Now you can list all envs no matter which provider they were using: 

```
kubectl get env
```

You should see `team-a` and `team-b` environments! :metal: