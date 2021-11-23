# Creating your first Knative Service

Once you have installed Knative Serving in your Kubernetes Cluster you can start creating Knative Services. 

In this short tutorial, we will be looking at what exactly happens when you create one, what is the default behavior and how it is different from creating a Kubernetes Deployment & Service pair. 

## Creating a Service and inspecting it

Go ahead and create a simple Knative Service by running the following command: 

```
kubectl apply -f email-knative-service.yaml
```

You can now list all your Knative Services by running: 
```
kubectl get ksvc
```

and you should see your newly created Knative Service:
```
NAME                  URL                                                        LATESTCREATED               LATESTREADY                 READY   REASON
email-service            http://email-service.default.X.X.X.X.sslip.io            email-service-00001            email-service-00001            True    

```

You can also describe the resource by running:

```
kubectl describe ksvc email-service
```

And you should see something like this: 

```
Name:         email-service
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
        Image:  salaboy/fmtok8s-email-rest:0.1.0
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
    URL:  http://email-service.default.svc.cluster.local
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
  Latest Created Revision Name:  email-service-00001
  Latest Ready Revision Name:    email-service-00001
  Observed Generation:           1
  Traffic:
    Latest Revision:  true
    Percent:          100
    Revision Name:    email-service-00001
  URL:                http://email-service.default.X.X.X.X.sslip.io
Events:
  Type    Reason   Age   From                Message
  ----    ------   ----  ----                -------
  Normal  Created  18m   service-controller  Created Configuration "email-service"
  Normal  Created  18m   service-controller  Created Route "email-service"
```

Check that the service is working using `curl` to the `info` endpoint:

```
curl http://email-service.default.X.X.X.X.sslip.io/info
```

This should return a json payload like the following: 

```
{"name":"Email Service (REST)","version":"v0.1.0","source":"https://github.com/salaboy/fmtok8s-email-rest/releases/tag/v0.1.0","podId":"","podNamepsace":"","podNodeName":""}
```

# Canary Releases with Knative (percentage-based split)

We want now to test a new version that our development teams have produced `salaboy/fmtok8s-email-rest:0.1.0-improved`. If we just change the Knative Service to use this new container image all the traffic will be automatically routed to it. Instead we can create some traffic routing rules to just send a percentage of the traffic to this new version until we are confident that things are working as expected. 

Edit or patch your Knative Service to update the traffic rules: 

```
kubectl edit ksvc email-service
```

Change the current `image` to `salaboy/fmtok8s-email-rest:0.1.0-improved` and modify traffic section to look like the following: 

```
  traffic: 
  - percent: 80
    revisionName: email-service-00001
  - latestRevision: true
    percent: 20
```

Alternatively you can apply the `email-knative-service-canary.yaml` resource: 

```
kubectl apply -f email-knative-service-canary.yaml
```

You should see the following output no matter which option did you used: 

```
service.serving.knative.dev/email-service configured
```

We have just configured 80% of the traffic to keep going to the stable version and only 20% to go to the newly produced version.

# A/B Testing (tags and header based routing)
