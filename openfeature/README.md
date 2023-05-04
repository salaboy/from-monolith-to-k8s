# Openfeature for Feature Flagging

On this tutorial we will be writing a small application that consume feature flags that are managed by OpenFeature, a CNCF initiative for feature flags backed up by multiple vendors.

## Installation

Based on the documentation that you can find in github
https://github.com/open-feature/open-feature-operator/blob/main/docs/installation.md


Prerequisites (CertManager):
```
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.10.1/cert-manager.yaml
kubectl wait --for=condition=Available=True deploy --all -n 'cert-manager'

```

Installing the OpenFeature K8s Operator with Helm: 

```
helm repo add openfeature https://open-feature.github.io/open-feature-operator/
helm repo update
helm upgrade --install openfeature openfeature/open-feature-operator
```

