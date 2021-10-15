# Creating your first Knative Service

Once you have installed Knative Serving in your Kubernetes Cluster you can start creating Knative Services. 

In this short tutorial, we will be looking at what exactly happens when you create one, what is the default behavior and how it is different from creating a Kubernetes Deployment & Service pair. 

## Creating a Service and inspecting it

Go ahead and create a simple Knative Service by running the following command: 

```
kubectl apply -f ui-knative-service.yaml
```

You can now list all your Knative Services by running: 
```
kubectl get ksvc
```

and you should see your newly created Knative Service:
```
NAME                  URL                                                        LATESTCREATED               LATESTREADY                 READY   REASON
ui-service            http://ui-service.default.X.X.X.X.sslip.io            ui-service-00001            ui-service-00001            True    

```

You can also describe the resource by running:

```
kubectl describe ksvc my-service
```

And you should see something like this: 

```
Name:         my-service
Namespace:    default
Labels:       <none>
Annotations:  serving.knative.dev/creator: salaboy@gmail.com
              serving.knative.dev/lastModifier: salaboy@gmail.com
API Version:  serving.knative.dev/v1
Kind:         Service
Metadata:
  Generation:          1
Spec:
  Template:
    Metadata:
      Creation Timestamp:  <nil>
    Spec:
      Container Concurrency:  0
      Containers:
        Image:  salaboy/fmtok8s-api-gateway:0.1.0
        Name:   user-container
        Readiness Probe:
          Success Threshold:  1
          Tcp Socket:
            Port:  0
        Resources:
      Enable Service Links:  false
      Timeout Seconds:       300
  Traffic:
    Latest Revision:  true
    Percent:          100
Status:
  Address:
    URL:  http://ui-service.default.svc.cluster.local
  Conditions:
    Last Transition Time:        2021-10-14T11:26:37Z
    Status:                      True
    Type:                        ConfigurationsReady
    Last Transition Time:        2021-10-14T11:26:38Z
    Status:                      True
    Type:                        Ready
    Last Transition Time:        2021-10-14T11:26:38Z
    Status:                      True
    Type:                        RoutesReady
  Latest Created Revision Name:  ui-service-00001
  Latest Ready Revision Name:    ui-service-00001
  Observed Generation:           1
  Traffic:
    Latest Revision:  true
    Percent:          100
    Revision Name:    ui-service-00001
  URL:                http://ui-service.default.34.65.77.156.sslip.io
Events:
  Type    Reason   Age   From                Message
  ----    ------   ----  ----                -------
  Normal  Created  18m   service-controller  Created Configuration "ui-service"
  Normal  Created  18m   service-controller  Created Route "ui-service"
```

