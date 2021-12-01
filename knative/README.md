# Deploying the Conference Platform with Knative

This document explains how to deploy the application Services as Knative Services. It also shows how you can do traffic splitting based on Headers. 

Then it goes to add Knative Eventing and enable the application to emit events to a Knative Broker which can be configured to send events to any other configured application. We will use Sockeye for example purposes. 


## Install Knative Serving

Follow the instructions here: 

https://knative.dev/docs/install/serving/install-serving-with-yaml/ 


Apply the following patch to support traffic splitting with Headers (explained here: https://knative.dev/docs/serving/feature-flags/#kubernetes-fieldref and  https://knative.dev/docs/serving/samples/tag-header-based-routing/) and Downward API (explained here: https://knative.dev/docs/serving/feature-flags/#kubernetes-fieldref):

```
kubectl patch cm config-features -n knative-serving -p '{"data":{"tag-header-based-routing":"Enabled", "kubernetes.podspec-fieldref": "Enabled"}}'
```

Note: you can use ModHeader to modify the request headers in your browser: https://chrome.google.com/webstore/detail/modheader/idgpnmonknjnojddfkpgkljpfnnfcklj?hl=en

### Install the application using Knative Services

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

You can list the Knative Services by running: 

```
kubectl get ksvc 
```
You should see something like: 
```
NAME                  URL                                                          LATESTCREATED               LATESTREADY                 READY   REASON
fmtok8s-agenda        http://fmtok8s-agenda.default.X.X.X.X.sslip.io        fmtok8s-agenda-00001        fmtok8s-agenda-00001        True    
fmtok8s-api-gateway   http://fmtok8s-api-gateway.default.X.X.X.X.sslip.io   fmtok8s-api-gateway-00001   fmtok8s-api-gateway-00001   True    
fmtok8s-c4p           http://fmtok8s-c4p.default.X.X.X.X.sslip.io           fmtok8s-c4p-00001           fmtok8s-c4p-00001           True    
fmtok8s-email         http://fmtok8s-email.default.X.X.X.X.sslip.io         fmtok8s-email-00001         fmtok8s-email-00001         True    

```
Where instead of `X`s you should see your public IP address.  
You can access the application by pointing your browser to: http://fmtok8s-api-gateway.default.X.X.X.X.sslip.io

You can send the following POST request to generate some talks proposals in the application: 
```
curl -X POST http://fmtok8s-api-gateway.default.X.X.X.X.sslip.io/api/test
```
Then go to the Back office section of the application and approve all the proposals. You should see them in the Main Site listed in different days. 

### Testing Traffic Splitting Using Percentages

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

```
This change will create a new revision, which we can use to split traffic. For doint that we need to add the following values into the `traffic` section:

```

```
  traffic:
  - latestRevision: false
    percent: 50
    revisionName: fmtok8s-api-gateway-00001
  - latestRevision: true
    percent: 50
```


### Testing Traffic Splitting Using Headers

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
image: salaboy/fmtok8s-api-gateway:0.1.0-debug
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
    revisionName: fmtok8s-api-gateway-00003
    tag: debug
  - latestRevision: true
    percent: 0
    tag: latest
```

With something like ModHeader for Chrome you can now specify the `debug` revision by setting the following header: 
`Knative-Serving-Tag` with value `debug`


## Installing Knative Eventing

Follow the instructions from here: 
https://knative.dev/docs/install/eventing/install-eventing-with-yaml/

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

## Configuring the project to use the Knative Eventing RabbitMQ Broker

To change the Broker implementation, and use the [https://github.com/knative-sandbox/eventing-rabbitmq/](https://github.com/knative-sandbox/eventing-rabbitmq/).

First we need to have the required CRDs for a RabbitMQ Operator to work:
- Install the RabbitMQ Cluster Operator
  kubectl apply -f https://github.com/rabbitmq/cluster-operator/releases/latest/download/cluster-operator.yml
- Install the Cert Manager required for the RabbitMQ Message Topology Operator, this because the TLS enabled admission webhooks needed for the Topology Operator to work properly
  kubectl apply -f https://github.com/jetstack/cert-manager/releases/latest/download/cert-manager.yaml
  kubectl wait --for=condition=Ready pods --all -n cert-manager
- Lastly, install the RabbitMQ Message Topology Operator
  kubectl apply -f https://github.com/rabbitmq/messaging-topology-operator/releases/latest/download/messaging-topology-operator-with-certmanager.yaml

Then, lets create a RabbitMQ Cluster:
```
kubectl create -f - <<EOF
  apiVersion: rabbitmq.com/v1beta1
  kind: RabbitmqCluster
  metadata:
    name: rabbitmq-cluster
    namespace: rabbitmq-resources
  spec:
    replicas: 1
EOF
```

Apply the RabbitMQ Broker CRD YAML:
```
kubectl apply -f https://github.com/knative-sandbox/eventing-rabbitmq/releases/download/knative-v1.0.0/rabbitmq-broker.yaml
```

Now lets create a RabbitMQ Broker:
```
kubectl create -f - <<EOF
  apiVersion: eventing.knative.dev/v1
  kind: Broker
  metadata:
    name: default
    namespace: rabbitmq-resources
    annotations:
      eventing.knative.dev/broker.class: RabbitMQBroker
  spec:
    config:
      apiVersion: rabbitmq.com/v1beta1
      kind: RabbitmqCluster
      name: rabbitmq-cluster
EOF
```

Now install the frontend (named app) and tickets charts using helm:
```
cat <<EOF | helm install app fmtok8s/fmtok8s-app --values=-
fmtok8s-api-gateway:
  knativeDeploy: true
  env:
    KNATIVE_ENABLED: "true"
    AGENDA_SERVICE: http://fmtok8s-agenda.default.svc.cluster.local
    C4P_SERVICE: http://fmtok8s-c4p.default.svc.cluster.local
    EMAIL_SERVICE: http://fmtok8s-email.default.svc.cluster.local
    K_SINK: http://default-broker-ingress.default.svc.cluster.local
    K_SINK_POST_FIX: "/broker, /"
    FEATURE_TICKETS_ENABLED: "true"

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

```
cat <<EOF | helm install tickets fmtok8s/fmtok8s-tickets --values=-
fmtok8s-tickets-service:
  knativeDeploy: true
  env:
    K_SINK: http://default-broker-ingress.default.svc.cluster.local
fmtok8s-payments-service:
  knativeDeploy: true
fmtok8s-queue-service:
  knativeDeploy: true
  env:
    K_SINK: http://default-broker-ingress.default.svc.cluster.local
EOF
```

## Debugging RabbitMQ

To debug RabbitMQ resources, fin the pod in the default namespace called
cluste-server-0, and port forward the port 15672:
kubectl port-forward cluster-server-0 15672:15672

Then find the RabbitMQ cluster default credentials, created when the Cluster yaml
was executed. This are located on the secret cluster-default-user in base64 encoding:
kubectl get secrets cluster-default-user -o json | jq -r '.data["default_user.conf"]' | base64 -d

Now go to http://localhost:15672/ and login with this credentials, here you have the
RabbitMQ Management UI were are the resources of RabbitMQ can be managed and monitored.

## Cleanup

To clean up this project resources use the next commands:
```
helm delete app tickets
```

And if you have the Knative Eventing RabbitMQ Broker implementation:
```
kubectl delete ns rabbitmq-resources
```
