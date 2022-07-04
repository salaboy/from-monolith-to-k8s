# Extending Kubernetes

This directory contains different Controllers implementations using different tools. 
All of them main objective is to monitor an instance of the Conference application (composed by 4 services) and report status back. 

## Kubebuilder


The [`kubebuilder`](https://book.kubebuilder.io/quick-start.html) directory contains a Kubernetes controller created with `kubebuilder` using Go. 

The instrutions for creating this controller from scratch are: 

- Install the [`kubebuilder` CLI](https://book.kubebuilder.io/quick-start.html)
- `mkdir conference-controller && cd conference-controller`
- `kubebuilder init --domain salaboy.com --repo github.com/salaboy/conference-controller`
- `kubebuilder create api --group conference --version v1 --kind Conference`
- `make manifests` 
- `make install` 
- `make run`
- `kubectl apply -f config/samples/` 

On top of the basic scaffolded project you can find the following logic implemented in this repository: 
- Get Conference Resource from the API Server
- if it exist:
  - Get all services that matches a label: "draft":"draft-app"
  - Check that each service in the app exist
  - If all services exist, for each service execute a request to understand if the service is working as expected (this should execute an operation)
  - If all services are operation mark the Conference.Status.Ready to true, if not mark it to false and emit a notification 
  - If a service with the name "fmtok8s-frontend" exist get its Status.LoadBalancer.Ingress[0].IP 


To the `ConferenceStatus struct` two new properties were added which require to clean up the CRD and re-run `make install`

```
Ready bool   `json:"ready"`
URL   string `json:"url"`
```

To show the Conference Status and URL annotations are needed to the `Conference struct`:

```
// +kubebuilder:printcolumn:name="READY",type="boolean",JSONPath=".status.ready"
// +kubebuilder:printcolumn:name="URL",type="string",JSONPath=".status.url"
```

Because this resource will not be modified, we need to requeue the reconcilation for a future point in time. In this case, we are setting a recurring period of 5 seconds by returning:
```
requeue := ctrl.Result{RequeueAfter: time.Second * 5}

return requeue, nil
```


