# Knative Eventing Tickets Sale Example

This tutorial install the Conference Platform application using Helm, but it also adds the services to implement the Ticket Sale flow. 

You can follow the Knative tutorial for installing the main Application Services: https://github.com/salaboy/from-monolith-to-k8s/tree/master/knative

Then to install the remaining services you can install the following Helm chart:

```
cat <<EOF | helm install app fmtok8s/fmtok8s-tickets --values=-
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


