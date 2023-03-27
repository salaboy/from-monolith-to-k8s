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

**Note**: you need to wait for the ArgoCD pods to be started, the first time you do this it will take more time, because it needs to fetch the container images from the internet.

You can access the user interface by pointing your browser to [http://localhost:8080](http://localhost:8080)

![](img/argocd-dashboard-login.png)

**Note**: by default the installation works using HTTP and not HTTPS, hence you need to accept the warning (hit the "Advanced Button" on Chrome) and proceed (Process to localhost unsafe). 


The user is `admin`, and to get the password for the ArgoCD Dashboard by running: 

```
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d; echo
```



# Setting up our application for the Staging Environment

On this short tutorial we will use a single namespace to represent our Staging Environment. We will store the configuration that we want to sync the Staging environment into a GitHub repository. 

First let's create a namespace for our Staging Environment:

```
kubectl create ns staging
```

Once you have ArgoCD installed you can access to the user interface to set up the project. 

![](img/argocd-dashboard-create-app.png)

Hit the "Create" button and use the following details to configure your project: 

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
kubectl port-forward svc/fmtok8s-frontend -n staging 8081:80
```

Then you can access the application pointing your browser to [http://localhost:8082](http://localhost:8082).


As usual, you can monitor the status of the pods and services using `kubectl`. To check if the application pods are ready you can run: 

```
kubectl get pods -n staging
```

To update version of configurations of your services, you can update the files located in the [Chart.yaml](staging/Chart.yaml) file or [values.yaml](staging/values.yaml) file located inside the [staging](staging/) directory. 




# Recap

