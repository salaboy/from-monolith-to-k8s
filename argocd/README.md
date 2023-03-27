# GitOps with ArgoCD

In this short tutorial we will set up the Conference application using individual Helm Charts for each Service. 
We will define our Staging environment using a Git repository. You can find here a [`staging` directory](staging/) which contains the definition of the 
charts that we want to deploy in our environment. 

We will be using ArgoCD to sync this directory that contains the definiton of what we want to deploy into a Kubernetes Cluster that represents our live Staging Environment for our 
development teams to test the application. 

# Prerequisites and installation

- Have a Kubernetes Cluster, we will use Kubernetes KinD in this tutorial
- Install ArgoCD in your cluster and install the CLI: https://argo-cd.readthedocs.io/en/stable/getting_started/
- You can fork/copy this repository http://github.com/salaboy/from-monolith-to-k8s/ as if you want to change the configuration for the application you will need to have write access to the repository. We will be using the directory `argocd/staging/`

Let's create a new KinD Cluster:

```
kind create cluster
```

Let's install ArgoCD in the cluster: 

```
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
```


You can access the ArgoCD User Interface by using port-forward, in a new terminal run:

```
kubectl port-forward svc/argocd-server -n argocd 8080:443
```

You can access the user interface by pointing your browser to [http://localhost:8080](http://localhost:8080)

**Note**: by default the installation works using HTTP and not HTTPS, hence you need to accept the warning and proceed. 


The user is `admin`, and to get the password for the ArgoCD Dashboard by running: 

```
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d; echo
```

As an exercise before getting things installed with ArgoCD you can install each individual chart using Helm to your cluster.

```
helm repo add fmtok8s https://salaboy.github.io/helm/
helm repo update
helm install c4p fmtok8s/fmtok8s-c4p-service 
helm install agenda fmtok8s/fmtok8s-agenda-service
helm install email fmtok8s/fmtok8s-email-service --version v0.0.1
helm install frontend fmtok8s/fmtok8s-frontend
```

You can access the app running using port-forwarding too: 

```
kubectl port-forward svc/fmtok8s-frontend 8081:80
```

So far you have: 
- ArgoCD Dashboard in port 8080
- Conference Application in port 8081

# Setting up our application for the Staging Environment

First let's create a namespace for our Staging Environment:

```
kubectl create ns staging
```

Once you have ArgoCD installed you can access to the user interface to set up the project. 

Hit the "Create Project" button and use the following details to configure your project: 

Here are the Create Application inputs that I've used: 
- Application Name: "staging-environment"
- Project: "default"
- Sync Policy: "Manual"
- Source Repository: "https://github.com/salaboy/from-monolith-to-k8s" (here you can point to your fork)
- Revision: "HEAD"
- Path: "argocd/staging/"
- Cluster: "https://kubernetes.default.svc" 
- Namespace: "staging"

And left the other values to their default ones. 




If you are running in a local environment, you can always access the application using `port-forward`:

```
kubectl port-forward svc/fmtok8s-frontend -n staging 8082:80
```

Then you can access the application pointing your browser to `http://localhost:8082`.



To get the external IP just list the services in the staging environment:

```
kubectl get svc -n staging fmtok8s-frontend 
NAME               TYPE           CLUSTER-IP    EXTERNAL-IP      PORT(S)        AGE
fmtok8s-frontend   LoadBalancer   10.116.0.45   X.X.X.X   80:30409/TCP   2m43s
```
The `EXTERNAL-IP` column should contain the IP that you can use to access the application.

Then you can access the application pointing your browser to `http://X.X.X.X`.

As usual, you can monitor the status of the pods and services using `kubectl`. To check if the application pods are ready you can run: 

```
kubectl get pods -n staging
```

