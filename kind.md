# Running the application with Kubernetes KIND

## Pre-Requisites
- `kubectl` installed. [For instructions check here](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- Helm 3 installed. [For instructions check here](https://helm.sh/docs/intro/install/)

## Creating a Kubernetes Cluster with KIND

Install [Kubernetes KIND](https://kind.sigs.k8s.io/docs/user/quick-start/#installation), this will allow you to create a new Kubernetes Cluster running on your laptop. 

```
cat <<EOF | kind create cluster --name dev --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
- role: worker
- role: worker
- role: worker
EOF
```

The cluster that you are creating will have 4 nodes, 3 workers and a control plane. 
This is to simulate a real cluster with a set of machines or virtual machines. 

![KIND Cluster creation](imgs/kind-cluster-creation.png)

In order to connect your `kubectl` CLI tool with this newly created you might need to run:

```
kubectl cluster-info --context kind-dev
```
![KIND Cluster Connect](imgs/kind-kubectl-connect.png)

Once you connected with the cluster you can start running commands against the cluster. For example you can check the cluster nodes by running:

```
kubectl get nodes -owide
```

![KIND Get Nodes oWide](imgs/kind-kubectl-get-nodes.png)

As you can see,  your Kubernetes Cluster is composed by 4 nodes and one of those is the control plane. 

Congrats your Cluster is up and running and you can connect with `kubectl`! 

## Installing the application

Now you are ready to install the application. 
You are going to install the application using Helm, a package manager for Kubernetes Applications. Helm allows you to install a complex Cloud-Native application and 3rd party software with a single command line. In order to install Helm Charts (packages/applications) you can add new repositories where your applications are stored. For java developers, these repositories are like Maven Central, Nexus or Artifactory. 

```
h repo add dev http://chartmuseum-jx.35.222.17.41.nip.io
h repo update
```

The previous two lines added a new repository to your Helm installation called `dev`, the second one fetched a file describing all the available packages and their versions in each repo that you have registered. 

Now that your Helm installation fecthed all the available packages description, let's install the application with the following line:

```
salaboy> helm install app dev/fmtok8s-app     
NAME: app
LAST DEPLOYED: Wed Dec 23 12:00:24 2020
NAMESPACE: default
STATUS: deployed
REVISION: 1
TEST SUITE: None
NOTES:
Cloud-Native Conference Platform V1

Chart Deployed: fmtok8s-app - 0.0.82
Release Name: app

```

Once the application is deployed, containers will need to be downloaded to your laptop in order to run, this can take a while. You can monitor the progress by listing all the pods running in your cluster, once again, using the `-owide` flag to get more information:

```
kubectl get pods -owide
```

![KIND Get Pods oWide](imgs/kind-kubectl-get-pods.png)

You need to pay attention to the `READY` and `STATUS` columns, where `1/1` in the `READY` column means that one replica of pod is correctly running and one was expected to be running. 

Notice that pods can be scheduled in different nodes in the `NODE` column, this is Kubernetes using the resources available in the cluster in an efficient way.



