# Openfeature for Feature Flagging

On this short tutorial we will enable our application to use Feature Flags by using Open Feature `flagd`. This enable each of the application services (including the Frontend) to consume and evaluate feature flags, in this case defined inside a Kubernertes configMap. 

## Installation

We will deploy a `flagd` Proxy which is in charge of mounting a Kubernetes `ConfigMap` that contains the Feature Flag definitions. 
Our application then can interact with this proxy to consume and evaluate the value of different feature flags.
Notice that this file also contains the `ConfigMap` with the feature flags.  

```
kubectl apply flagd.yaml
```

