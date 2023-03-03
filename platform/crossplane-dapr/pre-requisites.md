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