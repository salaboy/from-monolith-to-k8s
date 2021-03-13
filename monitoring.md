# Monitoring

When working with Cloud Native application understanding what is going on inside our services and being able to setup metrics and alerts is a must for our operation teams. 
In this short tutorial, we show the steps to setup to very popular tools (stacks of tools) which allows you to get more insights into how your applications are working internally inside the Kubernetes Cluster.

## Pre Requisites
You need to have a Kubernetes Cluster with the application deployed. Here are the steps summarized: 

1) Create a cluster with Kubernetes KIND

```
$ cat <<EOF | kind create cluster --name dev --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  kubeadmConfigPatches:
  - |
    kind: InitConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "ingress-ready=true"
  extraPortMappings:
  - containerPort: 80
    hostPort: 80
    protocol: TCP
  - containerPort: 443
    hostPort: 443
    protocol: TCP
- role: worker
- role: worker
- role: worker
EOF
```

2) Install Ingress Controller
```
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/static/provider/kind/deploy.yaml

```
3) Deploy the Application:

```
helm repo add dev http://chartmuseum-jx.34.67.22.199.nip.io
helm repo update
helm install app dev/fmtok8s-app
```



## Adding Prometheus to our application

We will install the Prometheus Stack which comes with Prometheus, Grafana and AlertsManager: 

```
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update

helm install my-prometheus prometheus-community/kube-prometheus-stack
```

Once the stack is up you can create a new Prometheus `ServiceMonitor` to start scrapping your services metrics:

```
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: fmtok8s-servicemonitor
  labels:
    release: my-prometheus
spec:
  selector:
    matchLabels:
      draft: draft-app
  namespaceSelector:
    matchNames:
      - default
  jobLabel: service-stats
  endpoints:
  - path: /actuator/prometheus
    targetPort: http
    interval: 15s
```

Notice that `metadata.labels[0]` -> `release: my-prometheus` needs to match with the name that you used for your Helm Release (`helm install my-prometheus`) 
Also notice that `spec.selector.matchLabels[0]` -> `draft: draft-app` This label needs to be added to all the services that you want to monitor (yes Kubernetes Services, not the Pods, not the Deployments) 

Also: `- path: /actuator/prometheus` is used because we are using Spring Booth actuators and the Prometheus Starter: 
```
<dependency>
   <groupId>io.micrometer</groupId>
   <artifactId>micrometer-registry-prometheus</artifactId>
   <scope>runtime</scope>
</dependency>
```

Important: `targetPort: http` makes reference to the containerPort inside the Pod, defined in the Deployment make sure the containerPort has a name, in this case `http`

To access Prometheus now you need to use port-forwarding:

```
kubectl port-forward svc/prometheus-operated 9090:9090
```

Check that your services are being scraped by going to **Status** -> **Service Discovery** you should see 
```
default/fmtok8s-servicemonitor/0 (4 / 25 active targets)
```

As this shows that **4** services are being scrapped. 


## Creating Grafana Dashboards

Grafana is already installed so let's access it using port-forward:

```
kubectl port-forward svc/my-prometheus-grafana 3000:80 
```

Point your browser to `http://localhost:3000`

The user and password are: **admin** and **prom-operator** respectively. You can find the password by looking at the secrets created: 

```
k edit secret my-prometheus-grafana
```

and then you can base64 decode the field: `admin-password`

Then you can click the **+** icon on the left of the screen and then Import to import a premade dashboard. Use the following file: [grafana-dashboard.json]()


## Adding Jaeger to our application


