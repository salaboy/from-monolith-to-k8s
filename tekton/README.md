# Pipelines



```
kubectl apply -f tekton/
```

Create Docker Hub secret: 

```
kubectl create secret docker-registry regcred --docker-server=https://index.docker.io/v1/ --docker-username=DOCKER_USERNAME --docker-password=DOCKER_PASSWORD --docker-email DOCKER_EMAIL
```

```
tkn pipeline start staging-environment-pipeline -w name=sources,volumeClaimTemplateFile=workspace-template.yaml -w name=dockerconfig,secret=regcred
```



# References
Why [JX uses Helmfile](https://jenkins-x.io/v3/develop/faq/general/#why-does-jenkins-x-use-helmfile-template)?

