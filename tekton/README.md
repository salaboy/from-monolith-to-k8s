# Pipelines

This document explains how to run a Service and an Environment Pipeline using Tekton. 

## Installing Tekton

1. Install Tekton Pipeline:
```
kubectl apply -f https://storage.googleapis.com/tekton-releases/pipeline/previous/v0.25.0/release.yaml
```
2. Install Tekton Triggers:
```
kubectl apply -f https://storage.googleapis.com/tekton-releases/triggers/previous/v0.14.2/release.yaml
kubectl apply -f https://storage.googleapis.com/tekton-releases/triggers/previous/v0.14.2/interceptors.yaml
```
3. Install Tekton Dashboard (optional):
```
kubectl apply -f https://github.com/tektoncd/dashboard/releases/download/v0.17.0/tekton-dashboard-release.yaml
```

Install Tekton `tkn` CLI tool: https://github.com/tektoncd/cli

### Configure Tekton Pipeline

The Tekton pipeline uses some alpha features (Tekton Bundles) which needs to
be enabled in the config. The `feature-flags` config map in the `tekton-pipelines` namespace
should look like:

```yaml
apiVersion: v1
data:
  enable-api-fields: alpha # <------- 
  disable-affinity-assistant: "false"
  disable-creds-init: "false"
  disable-home-env-overwrite: "true"
  disable-working-directory-overwrite: "true"
  enable-custom-tasks: "false"
  enable-tekton-oci-bundles: "false"
  require-git-ssh-secret-known-hosts: "false"
  running-in-environment-with-injected-sidecars: "true"
kind: ConfigMap
(...)
```

## RBAC

If the pipeline is going to push docker images to DockerHub you need the following steps: 


Create Docker Hub secret: 

```
kubectl create secret docker-registry regcred --docker-server=https://index.docker.io/v1/ --docker-username=DOCKER_USERNAME --docker-password=DOCKER_PASSWORD --docker-email DOCKER_EMAIL
```

Then apply all the RBAC configurations and the pipelines: 


```
kubectl apply -f tekton/
```

## Service Pipeline

The Service Pipeline definition described in `resources/service-pipeline.yaml` implements the following tasks:

<Diagram>

You can start this Service Pipeline by running the following command:

```
tkn pipeline start api-gateway-service-pipeline -s dockerconfig -w name=sources,volumeClaimTemplateFile=workspace-template.yaml -w name=dockerconfig,secret=regcred -w name=maven-settings,emptyDir=
```



# References
Why [JX uses Helmfile](https://jenkins-x.io/v3/develop/faq/general/#why-does-jenkins-x-use-helmfile-template)?

