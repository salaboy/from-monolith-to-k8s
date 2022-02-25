# Super fast getting started guide with `func` and Go

```
mkdir uppercase
cd uppercase/
```

```
func create -l go -t uppercase --repository https://github.com/salaboy/func-templates
```

```
func build 
```

```
func run
```

```
curl -v -X POST -d '{"input": "salaboy"}' \
  -H'Content-type: application/json' \
  -H'Ce-id: 1' \
  -H'Ce-source: cloud-event-example' \
  -H'Ce-subject: Convert to UpperCase' \
  -H'Ce-type: UppercaseRequestedEvent' \
  -H'Ce-specversion: 1.0' \
  http://localhost:8080/
```

```
func deploy
```

```
mkdir improve
cd improve/
```

```
func create -l go -t improve --repository https://github.com/salaboy/func-templates
```

```
func build 
```

```
func run
```

```
curl -v -X POST -d '{"input": "salaboy", "output": "SALABOY", "operation": "uppercase"}' \
  -H'Content-type: application/json' \
  -H'Ce-id: 1' \
  -H'Ce-source: cloud-event-example' \
  -H'Ce-subject: Convert to UpperCase' \
  -H'Ce-type: UpperCasedEvent' \
  -H'Ce-specversion: 1.0' \
  http://localhost:8080/
```


```
func deploy
```

```
kubectl logs -f improve<TAB> user-container
```

```
kubectl create -f - <<EOF
apiVersion: eventing.knative.dev/v1
kind: Broker
metadata:
 name: default
 namespace: default
EOF
```

```
kubectl create -f - <<EOF
apiVersion: eventing.knative.dev/v1
kind: Trigger
metadata:
  name: uppercase-trigger
  namespace: default
spec:
  broker: default
  filter:
    attributes:
      type: UppercaseRequestedEvent
  subscriber:
    ref:
      apiVersion: serving.knative.dev/v1
      kind: Service
      name: uppercase
--- 

apiVersion: eventing.knative.dev/v1
kind: Trigger
metadata:
  name: improve-trigger
  namespace: default
spec:
  broker: default
  filter:
    attributes:
      type: UpperCasedEvent
  subscriber:
    ref:
      apiVersion: serving.knative.dev/v1
      kind: Service
      name: improve

EOF
```

```
kubectl port-forward svc/broker-ingress -n knative-eventing 8080:80
```

```
curl -v "http://localhost:8080/default/default" \
-H "Content-Type:application/json" \
-H "Ce-Id:1" \
-H "Ce-Subject:Uppercase" \
-H "Ce-Source:cloud-event-example" \
-H "Ce-Type: UppercaseRequestedEvent" \
-H "Ce-Specversion:1.0" \
-d "{\"input\": \"salaboy\"}"

```

```
kubectl logs -f improve user-container
```
