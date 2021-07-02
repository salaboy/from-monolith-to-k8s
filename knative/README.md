# Deploying the Conference Platform with Knative

This document explains how to deploy the application Services as Knative Services. It also shows how you can do traffic splitting based on Headers. 

Then it goes to add Knative Eventing and enable the application to emit events to a Knative Broker which can be configured to send events to any other configured application. We will use Sockeye for example purposes. 


## Install Knative Serving

Follow the instructions here: 

https://knative.dev/docs/install/install-serving-with-yaml/

Apply the following patch to support traffic splitting with Headers (explained here: https://knative.dev/docs/serving/samples/tag-header-based-routing/):

```
kubectl patch cm config-features -n knative-serving -p '{"data":{"tag-header-based-routing":"Enabled"}}'
```

Note: you can use ModHeader to modify the request headers in your browser: https://chrome.google.com/webstore/detail/modheader/idgpnmonknjnojddfkpgkljpfnnfcklj?hl=en

### Install the application using Knative Services

```
cat <<EOF | helm install app fmtok8s/fmtok8s-app --values=-
fmtok8s-api-gateway:
  knativeDeploy: true
  env:
    AGENDA_SERVICE: http://fmtok8s-agenda.default.svc.cluster.local
    C4P_SERVICE: http://fmtok8s-c4p.default.svc.cluster.local
    EMAIL_SERVICE: http://fmtok8s-email.default.svc.cluster.local

fmtok8s-agenda-rest:
  knativeDeploy: true
fmtok8s-c4p-rest:
  knativeDeploy: true
  env:
    AGENDA_SERVICE: http://fmtok8s-agenda.default.svc.cluster.local
    EMAIL_SERVICE: http://fmtok8s-email.default.svc.cluster.local
fmtok8s-email-rest:
  knativeDeploy: true
EOF
```

### Testing Traffic Splitting by Headers or Percentage

You can edit the Knative Service (ksvc) of the API Gateway and create a new revision by changing the docker image that the service is using: 

```
kubectl edit ksvc fmtok8s-api-gateway
```

Then modify the `image` name with the following value: 

From:
```
image: salaboy/fmtok8s-api-gateway:0.1.0
```
To:
```
image: salaboy/fmtok8s-api-gateway:0.1.0-color
```
This change will create a new revision, which we can use to split traffic. For doint that we need to add the following values into the `traffic` section:

```
  traffic:
  - latestRevision: false
    percent: 100
    revisionName: fmtok8s-api-gateway-00001
    tag: current
  - latestRevision: false
    percent: 0
    revisionName: fmtok8s-api-gateway-00002
    tag: candidate
  - latestRevision: true
    percent: 0
    tag: latest
```

## Installing Knative Eventing

Follow the instructions from here: 
https://knative.dev/docs/install/install-eventing-with-yaml/

Then create a Broker:

```
kubectl create -f - <<EOF
apiVersion: eventing.knative.dev/v1
kind: broker
metadata:
 name: default
 namespace: default
EOF
```

Then you can enable emitting events into the services by adding two environment variables.

```
kubectl edit ksvc fmtok8s-c4p
```
Then edit the `env` section and add: 
```
        - name: EVENTS_ENABLED
          value: "true"
        - name: K_SINK
          value: http://broker-ingress.knative-eventing.svc.cluster.local/default/default  
```

Now you can deploy Sockeye to monitor CloudEvents: 

https://github.com/n3wscott/sockeye

```
kubectl apply -f https://github.com/n3wscott/sockeye/releases/download/v0.7.0/release.yaml

```

Finally, you just need to create a trigger (subscription) to connect the consumer, in this case Sockeye to the Broker. Notice that the producer only knows where the Broker is. 

```
apiVersion: eventing.knative.dev/v1
kind: Trigger
metadata:
  name: c4p-trigger
  namespace: default
spec:
  broker: default
  subscriber:
    uri: http://sockeye.default.svc.cluster.local
```