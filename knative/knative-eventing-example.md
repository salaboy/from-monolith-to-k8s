# Knative Eventing Tickets Sale Example


This tutorial install the Conference Platform application using Helm, but it also adds the services to implement the Ticket Sale flow. 

## Pre Requisites
- Install Knative Serving and Knative Eventing
- Install the Conference Platform App using Helm, You can follow the Knative tutorial for installing the main Application Services: https://github.com/salaboy/from-monolith-to-k8s/tree/master/knative

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

You need to create a broker and have sockeye if you want to see the events flowing. 


## Installing the Tickets Queue Services

Then to install the remaining services you can install the following Helm chart:

```
cat <<EOF | helm install tickets-app fmtok8s/fmtok8s-tickets --values=-
fmtok8s-tickets-service:
  knativeDeploy: true
  env:
    KNATIVE_ENABLED: "true"
fmtok8s-payments-service:
  knativeDeploy: true
fmtok8s-queue-service:
  knativeDeploy: true

EOF
```


To enable the Tickets section in the application you should update the Knative Service called `fmtok8s-api-gateway` to have the following environment variable set: 

```
- name: TICKETS_ENABLED
  value: "true"
```
