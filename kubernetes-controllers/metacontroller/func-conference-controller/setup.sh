#!/usr/bin/env bash

##### This script will guide you through installing the prerequisites and running an example
##### It includes:
##### - validating required binaries are installed
##### - creating a local kind cluster with Knative Serving using Knative Quick Start
##### - installing Metacontroller
##### - deploying controller function
##### - registering function with Metacontroller
##### - creating Conference CRD
##### - deploying a Conference Application
##### - creating a Conference CRD instance
##### - validating the function was launched and that it completed successfully

# Get latest metacontroller-helm version here: https://github.com/metacontroller/metacontroller/pkgs/container/metacontroller-helm
METACONTROLLER_HELM_CHART_VERSION=4.5.4
IMAGE=docker.io/ciberkleid736/func-conference-controller:latest
# Choose builder for AMD64 (paketo) or ARM64 (dmikusa):
#BUILDER=paketobuildpacks/builder:base
BUILDER=dmikusa2pivotal/builder:focal

### Verify required binaries are installed

if ! command -v helm &> /dev/null; then echo "helm could not be found"; exit; fi
if ! command -v func &> /dev/null; then echo "func could not be found"; exit; fi
if ! command -v kn &> /dev/null; then echo "kn could not be found"; exit; fi
if ! command -v kubectl &> /dev/null; then echo "kubectl could not be found"; exit; fi
if ! command -v yq &> /dev/null; then echo "yq could not be found"; exit; fi

# Install Knative Serving per instructions at https://knative.dev
#   Can also use QuickStart:
#   https://knative.dev/docs/getting-started/quickstart-install/#install-the-knative-cli
if [[ $(ktx kind-fmtok8s-metacontroller) == "" ]]
then
  kn quickstart kind --name fmtok8s-metacontroller --install-serving
else
  echo "Using existing cluster (ktx kind-fmtok8s-metacontroller)"
fi

# Fetch metacontroller chart and install metacontroller
HELM_EXPERIMENTAL_OCI=1 helm pull oci://ghcr.io/metacontroller/metacontroller-helm --version=v${METACONTROLLER_HELM_CHART_VERSION}
kubectl create ns metacontroller
helm install metacontroller metacontroller-helm-v${METACONTROLLER_HELM_CHART_VERSION}.tgz --namespace metacontroller
rm metacontroller-helm-v${METACONTROLLER_HELM_CHART_VERSION}.tgz

# Update and apply func.yaml
yq -i 'del(.image)' func.yaml
yq -i 'del(.builder)' func.yaml
yq -i 'del(.builders)' func.yaml
yq -i 'del(.buildpacks)' func.yaml
echo "image: $IMAGE" >> func.yaml
echo -e "builderImages:\n  pack: $BUILDER" >> func.yaml

# Deploy controller function
# This function creates Deployments that run tests for Conference Applications 
func deploy -v

# Validate that the function was deployed
func list
# You can also run: kubectl get kservice

# Register the function with Metacontroller so that Metacontroller will 
# forward any events that are sent to the API Server for a Conference type resource
kubectl apply -f config/controller.yaml

# Validate that the controller was registered by confirming that a 
# CompositeController resource was created
kubectl get CompositeController

# Register "Conference" CRD with the Kubernetes API
kubectl apply -f config/crd.yaml

# Validate that the CRD was created
kubectl get crds conferences.metacontroller.conference.salaboy.com

# At this pont, we have the mechanics to support the following:
# - User creates, updates or deletes a Conference resource (instance of Conference CRD)
# - Metacontroller launches controller function to handle event

# What does the controller function do?
# Notice that the Conference CRD spec contains two values: a namespace and a boolean called productionTestsEnabled
# This function expects to find 4 apps in the given namespace (part of a "Conference Application," deployed separately).
# It checks the status of the 4 apps to make sure they are all "ready"
# If all apps are ready and productionTestsEnabled is true, then the function
# creates a Deployment to run tests against the apps.

# Of course, this means that for the function to be effective, we need to deploy
# at least one instance of a Conference Application for it to act on.

# In a separate namespace, deploy an instance of the Conference Application
# A Conference Application is defined by the helm chart below and 
# comprises 4 apps, postgres db, and redis
helm repo add fmtok8s https://salaboy.github.io/helm/
helm install conference fmtok8s/fmtok8s-conference-chart --namespace conf-jbcnconf --create-namespace

# Wait until all 4 conference apps are ready (status=running)
kubectl get pods -n conf-jbcnconf

# Next, create Conference resource to trigger the function to be launched and
# check on the Conferene Application
# Before you do that, open a separate terminal window and run the following command
# so that you can watch the function get launched by metacontroller when you
# create the Conference resource
kubectl get pods -w

# Next, create the conference resource
kubectl apply -f config/conference-jbcnconf.yaml

# Keep an eye on the terminal window where you are watching pods. You should see the
# function being launched. It might look something like this:

# Did it work? You should see a new deployment that 
# was created by the function in order to test the conference apps
kubectl get deployments

##### Try it again for a second conference!
# 1. Deploy the Conference Application
# 2. Create the Conference resource
# 3. Check for a Deployment
helm install conference fmtok8s/fmtok8s-conference-chart --namespace conf-springone --create-namespace
kubectl apply -f config/conference-springone.yaml
kubectl get deployments