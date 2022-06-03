# GitOps with ArgoCD

In this short tutorial we will set up the Conference application using individual Helm Charts for each Service. 
We will define our Staging environment using a Git repository. You can find here a `staging` directory which contains the definition of the 
charts that we want to deploy in our environment. 

We will be using ArgoCD to sync this directory that contains the definiton of what we want to deploy into a Kubernetes Cluster that represents our live Staging Environment for our 
development teams to test the application. 

# Prerequisites

- Install ArgoCD in your cluster and install the CLI: https://argo-cd.readthedocs.io/en/stable/getting_started/
- 


# Setting up our application for the Staging Environment

