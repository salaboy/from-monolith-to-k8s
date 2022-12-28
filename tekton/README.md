# Pipelines

This short tutorial shows how to use Tekton to define and run a Service and an Environment Pipeline. 

[Tekton](https://tekton.dev) is a non-opinionated Pipeline Engine built for the Cloud (specifically for Kubernetes). You can build any kind of pipelines that you want as the engine doesn't impose any restrictions on the kind of Tasks that it can execute. This makes it perfect for building Service Pipelines where you might need to have special requirements that cannot be met by a managed service.  

The Service Pipeline for this example is configured to build the [Conference Application Frontend](https://github.com/salaboy/fmtok8s-frontend) but as you can see in the [Service Pipeline definition](resources/service-pipeline.yaml) you can parameterize the pipeline run to build other services. 

The [Environment Pipeline definition](resources/environment-pipeline.yaml) shows a simple example on how you can use Helm to sync the contents of a repository to a namespace in a Kubernetes Cluster. While this is doable with Tekton, there are other more specialized tools like [ArgoCD](https://argo-cd.readthedocs.io/en/stable/) which do a more complete set of tasks on the continuous deployment space by applying a GitOps approach. You can find a [tutorial with ArgoCD here](../argocd/README.md).



## Installing Tekton

Follow the next steps in order to install and setup Tekton in your Kubernetes Cluster.

1. **Install Tekton Pipelines**

```
  kubectl apply -f https://storage.googleapis.com/tekton-releases/pipeline/previous/v0.42.0/release.yaml
```
2. **Install Tekton Triggers**

```
kubectl apply -f https://storage.googleapis.com/tekton-releases/triggers/previous/v0.22.0/release.yaml
kubectl apply -f https://storage.googleapis.com/tekton-releases/triggers/previous/v0.22.0/interceptors.yaml
```
3. **Install Tekton Dashboard (optional)**

```
kubectl apply -f https://github.com/tektoncd/dashboard/releases/download/v0.31.0/tekton-dashboard-release.yaml
```
You can access the dashboard by port-forwarding using `kubectl`:

```
kubectl port-forward svc/tekton-dashboard  -n tekton-pipelines 9097:9097
```

![Tekton Dashboard](tekton-dashboard.png)

Then you can access pointing your browser to [http://localhost:9097](http://localhost:9097)


4. **Install Tekton CLI (optional)**:

You can also install [Tekton `tkn` CLI tool](https://github.com/tektoncd/cli)

### Configure Tekton Pipeline

The Tekton pipeline definition uses Tekton Bundles which needs to
be enabled in the config. The `feature-flags` config map in the `tekton-pipelines` namespace
should look like:

```
kubectl edit cm -n tekton-pipelines feature-flags
```

```yaml
apiVersion: v1
data:
  enable-api-fields: stable 
  disable-affinity-assistant: "false"
  disable-creds-init: "false"
  disable-home-env-overwrite: "true"
  disable-working-directory-overwrite: "true"
  enable-custom-tasks: "false"
  enable-tekton-oci-bundles: "true" # <------- 
  require-git-ssh-secret-known-hosts: "false"
  running-in-environment-with-injected-sidecars: "true"
kind: ConfigMap
(...)
```

Check the [official documentation](https://github.com/tektoncd/pipeline/blob/release-v0.18.x/docs/install.md#customizing-the-pipelines-controller-behavior) for more information.

## RBAC

If the pipeline is going to push docker images to DockerHub you need the following steps: 

Create Docker Hub secret: 

```
kubectl create secret docker-registry regcred --docker-server=https://index.docker.io/v1/ --docker-username=DOCKER_USERNAME --docker-password=DOCKER_PASSWORD --docker-email DOCKER_EMAIL
```

To create this, in my Mac OSX laptop I need to access the `Keychain Access` app and then look at my `Docker Credentials`. This are generated when doing `docker login`. The DOCKER_PASSWORD is this hash, instead of my textual password for Docker Hub.

Then apply all the RBAC configurations and the pipelines: 

```
kubectl apply -f tekton/resources/
```

## Service Pipeline

The Service Pipeline definition described in [`resources/service-pipeline.yaml`](resources/service-pipeline.yaml) implements the following tasks:

![Service Pipeline](service-pipeline.png)

You can start this Service Pipeline by running the following command:

```
tkn pipeline start frontend-service-pipeline -s dockerconfig -w name=sources,volumeClaimTemplateFile=workspace-template.yaml -w name=dockerconfig,secret=regcred -w name=maven-settings,emptyDir=
```

## Environment Pipeline

The environment pipeline definition described in [`resources/envionment-pipeline.yaml`](resources/envionment-pipeline.yaml) implements the following tasks:

![Environment Pipeline](environment-pipeline.png)


You can start this Environment Pipeline by running the following command:

```
tkn pipeline start staging-environment-pipeline -w name=sources,volumeClaimTemplateFile=workspace-template.yaml -s gitops
```

The environment pipeline is using [`helmfile`](https://github.com/roboll/helmfile) to describe the stating environment. 


As mentioned before, while this is doable with Tekton, there are other more specialized tools like [ArgoCD](https://argo-cd.readthedocs.io/en/stable/) which do a more complete set of tasks on the continuous deployment space by applying a GitOps approach. You can find a [tutorial with ArgoCD here](../argocd/README.md).

# References
- Why [JX uses Helmfile](https://jenkins-x.io/v3/develop/faq/general/#why-does-jenkins-x-use-helmfile-template)?

