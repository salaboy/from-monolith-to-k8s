# Using Helm Provider from Crossplane

This example shows how you can use the [Helm Provider](https://doc.crds.dev/github.com/crossplane-contrib/provider-helm) to install the Conference Platform in separate namespaces using a declarative approach. 

## Installing and configuring the `provider-helm`
You first need to install and configure [`provider-helm`](https://doc.crds.dev/github.com/crossplane-contrib/provider-helm) into Crossplane. Run the following command to install the provider into your Crossplane instalation: 

```
kubectl apply -f helm-provider.yaml
kubectl apply -f helm-provider-config.yaml
```

You need to grant the provider the right access to deploy charts on your behalf: 

```
# Make sure provider-helm has enough permissions to install your chart into cluster
#
# You can give admin permissions by running:
# SA=$(kubectl -n crossplane-system get sa -o name | grep provider-helm | sed -e 's|serviceaccount\/|crossplane-system:|g')
# kubectl create clusterrolebinding provider-helm-admin-binding --clusterrole cluster-admin --serviceaccount="${SA}"
```

## Installing charts in declarative way
Once the provider is ready and configured, we can create a new `Helm Release` (an full instance of the Conference Platform) in different namespaces inside the same Kubernetes Cluster. 

Check the `conference-platform-release.yaml` file that describe which Helm Chart will be installed. 

```
apiVersion: helm.crossplane.io/v1beta1
kind: Release
metadata:
  name: conference-platform
spec:
  forProvider:
    chart:
      name: fmtok8s-app
      repository: https://salaboy.github.io/helm/
      version: 0.1.0
    namespace: conference-customer-a
    values:
      service:
        type: ClusterIP
  providerConfigRef:
    name: helm-provider
```

You can set all the parameters that the chart accept inside the Release resource (by setting the properties `spec.forProvider.set` which is an array of name/value pairs)

Then create a Helm Release for the Helm Provider to install:
```
kubectl apply -f conference-platform-release.yaml
```

You can query to see all the releases installed in the cluster with: 

```
salaboy> kubectl get release
NAME                  CHART         VERSION   SYNCED   READY   STATE      REVISION   DESCRIPTION        AGE
conference-platform   fmtok8s-app   0.1.0     True     True    deployed   1          Install complete   74s
```


And also check that a new namespace called `conference-customer-a` was created and contain the application pods:

```
salaboy> kubectl get pods -n conference-customer-a
NAME                                                       READY   STATUS      RESTARTS   AGE
conference-platform-fmtok8s-agenda-rest-64cccb7764-cnx25   1/1     Running     0          41s
conference-platform-fmtok8s-api-gateway-5f956fd899-xmcjc   1/1     Running     0          41s
conference-platform-fmtok8s-c4p-rest-578bc5db6c-7x6tm      1/1     Running     0          41s
conference-platform-fmtok8s-email-rest-6bdfb6558b-n57h2    1/1     Running     0          41s
```

If you made it this far, you can now install and create more advanced compositions mixing Helm Charts and Cloud Provider specific resources. 



