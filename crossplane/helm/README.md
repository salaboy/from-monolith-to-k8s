# Using Helm Provider from Crossplane

This example shows how you can use the Helm Provider to install the Conference Platform in separate namespaces using a declarative approach. 

We can create a new Helm Release for each instalation in different namespaces inside the same Kubernetes Cluster. 

Check the `release.yaml` file that describe which Helm Chart will be installed. 


Check provider config to enable the provider to install charts:

```

```

```
https://github.com/crossplane-contrib/provider-helm/blob/master/examples/provider-config/provider-config-incluster.yaml

```

Then create a Helm Release for the Helm Provider to install:
```
kubectl apply -f release.yaml
```
