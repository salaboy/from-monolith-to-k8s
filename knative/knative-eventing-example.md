# Knative Eventing Tickets Sale Example


This tutorial install the Conference Platform application using Helm, but it also adds the services to implement the Tickets Selling flow. 

## Pre Requisites
- Install [Knative Serving](https://knative.dev/docs/install/serving/install-serving-with-yaml/) and [Knative Eventing](https://knative.dev/docs/install/eventing/install-eventing-with-yaml/).
- Install the Conference Platform App using Helm and setting the `knativeDeploy` variable to `true`
- Create a Knative Eventing Broker, install SockEye and create a trigger to see all the events

### Creating a Knative Eventing Broker

```
kubectl create -f - <<EOF
apiVersion: eventing.knative.dev/v1
kind: broker
metadata:
 name: default
 namespace: default
EOF
```

### Installing the base Conference Platform using Knative Resources
```
cat <<EOF | helm install app fmtok8s/fmtok8s-app --values=-
fmtok8s-api-gateway:
  knativeDeploy: true
  env:
    KNATIVE_ENABLED: "true"
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
### Installing Sockeye for monitoring events

```
kubectl apply -f https://github.com/n3wscott/sockeye/releases/download/v0.7.0/release.yaml
```

### Creating a trigger to see all the events going to the broker

```
apiVersion: eventing.knative.dev/v1
kind: Trigger
metadata:
  name: wildcard-trigger
  namespace: default
spec:
  broker: default
  subscriber:
    uri: http://sockeye.default.svc.cluster.local
```

## Installing the Tickets Queue Services

Then to install the remaining services you can install the following Helm chart:

```
cat <<EOF | helm install tickets-app fmtok8s/fmtok8s-tickets --values=-
fmtok8s-tickets-service:
  knativeDeploy: true
fmtok8s-payments-service:
  knativeDeploy: true
fmtok8s-queue-service:
  knativeDeploy: true

EOF
```


To enable the Tickets section in the application you should update the Knative Service called `fmtok8s-api-gateway` to have the following environment variable set: 

```
- name: FEATURE_TICKETS_ENABLED
  value: "true"
```
