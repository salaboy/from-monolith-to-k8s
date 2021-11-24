# Installing Crossplane in GCP

This tutorial expands on the official documentation by also enabling Redis Services inside Google Cloud, so Crossplane can provision Redis instances for us. 


Install [Crossplane Self-Hosted following the instructions in the official site](https://crossplane.io/docs/v1.5/getting-started/install-configure.html). I've selected Self-Hosted and then Helm 3 for the installation in GKE. 
You should see something like this: 

```
salaboy> helm install crossplane --namespace crossplane-system crossplane-stable/crossplane
NAME: crossplane
LAST DEPLOYED: Tue Nov 23 13:17:50 2021
NAMESPACE: crossplane-system
STATUS: deployed
REVISION: 1
TEST SUITE: None
NOTES:
Release: crossplane

Chart Name: crossplane
Chart Description: Crossplane is an open source Kubernetes add-on that enables platform teams to assemble infrastructure from multiple vendors, and expose higher level self-service APIs for application teams to consume.
Chart Version: 1.5.0
Chart Application Version: 1.5.0

Kube Version: v1.21.5-gke.1302

```

**Note**: Make sure that you also install the Crossplane `kubectl` plugin. 

Proceed with the Getting Started Packages, because I am on GCP I will choose the GCP Configuration package.

```
salaboy> kubectl crossplane install configuration registry.upbound.io/xp/getting-started-with-gcp:v1.5.0
configuration.pkg.crossplane.io/xp-getting-started-with-gcp created
```

## Giving Crossplane rights to create resources on our behalf


**Note**: The following enable Redis API, make sure you run the commands in the right order so the SA ends up with these Service and Role too. 

```
export SERVICE="redis.googleapis.com"
gcloud services enable $SERVICE --project $PROJECT_ID

export ROLE="roles/redis.admin"
gcloud projects add-iam-policy-binding --role="$ROLE" $PROJECT_ID --member "serviceAccount:$SA"

```
Then you can create the Service Account Key file, as instructed in the docs: 

```
# create service account keyfile
gcloud iam service-accounts keys create creds.json --project $PROJECT_ID --iam-account $SA
```

**Note**: Make sure to not store this key file in Github, as this contains the credentails to use GCP on your behalf

Finally, create a Kubernetes Secret using this credentials (`creds.json`):

```
kubectl create secret generic gcp-creds -n crossplane-system --from-file=creds=./creds.json

```

Once the credentials are inside Kubernetes, we can configure our provider, in this case `provider-gcp`

```
# replace this with your own gcp project id
PROJECT_ID=my-project
echo "apiVersion: gcp.crossplane.io/v1beta1
kind: ProviderConfig
metadata:
  name: default
spec:
  projectID: ${PROJECT_ID}
  credentials:
    source: Secret
    secretRef:
      namespace: crossplane-system
      name: gcp-creds
      key: creds" | kubectl apply -f -
```

If you manage to finish these steps, your Crossplane instalation is ready to provision PostgreSQL and Redis instances for you. As you might have noticed, in the documentation the steps for AWS and Azure are pretty similar to GCP and it also important to note that they are not exclusive, you can configure as many providers as you want in your Crossplane installation. 

Check the official documentation on [how to provision infrastructure by applying GCP specific Kubernetes resoureces](https://crossplane.io/docs/v1.5/getting-started/provision-infrastructure.html). 

Also check the reference documentation about the [Resources that the GCP provider installs for you](https://doc.crds.dev/github.com/crossplane/provider-gcp). Notice that Crossplane will be able to provision only the resources that you enabled in the previous installation steps (the GCP APIs that we enabled and the Service Account that we included in the `creds.json` file contains which kind of rights, over which resources Crossplane can use/create). 


