# Dapr-enabled Environment Promise

This promise set up a Redis Database and a Dapr Statestore component for your developers to use right away.

After installing the promise in your cluster you can request new environments by sending requests like this: 

```
apiVersion: marketplace.kratix.io/v1alpha1
kind: env
metadata:
  name: my-dev-env
  namespace: default
spec:
  deployApp: true
  database:   
    enabled: true
    statestoreName: "statestore"
```

