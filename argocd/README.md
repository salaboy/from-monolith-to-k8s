# GitOps with ArgoCD

In this short tutorial we will set up the Conference application using individual Helm Charts for each Service. 
We will define our Staging environment using a Git repository. You can find here a `staging` directory which contains the definition of the 
charts that we want to deploy in our environment. 

We will be using ArgoCD to sync this directory that contains the definiton of what we want to deploy into a Kubernetes Cluster that represents our live Staging Environment for our 
development teams to test the application. 

# Prerequisites

- Install ArgoCD in your cluster and install the CLI: https://argo-cd.readthedocs.io/en/stable/getting_started/
- You can fork/copy this repository http://github.com/salaboy/from-monolith-to-k8s/ as if you want to change the configuration for the application you will need to have write access to the repository. We will be using the directory `argocd/staging/`

As an exercise before getting things installed with ArgoCD you can install each individual chart using Helm to your cluster.

```
helm repo add fmtok8s https://salaboy.github.io/helm/
helm repo update
helm install c4p fmtok8s/fmtok8s-c4p-service 
helm install agenda fmtok8s/fmtok8s-agenda-service
helm install email fmtok8s/fmtok8s-email-service --version v0.0.1
helm install frontend fmtok8s/fmtok8s-frontend
```


# Setting up our application for the Staging Environment

Once you have ArgoCD installed you can access to the user interface to set up the project. 
