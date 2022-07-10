# Extending Kubernetes

This directory contains different Controllers implementations using different tools. 
All of them main objective is to monitor an instance of the Conference application (composed by 4 services) and report status back. 

## Java Operator SDK

The Java Operator SDK allows us to create our Kubernetes Controller using the Fabric8.io Kubernetes APIs and their own abstractions. This doesn't use the Controller Runtime approach for interacting with the Kubernetes API Server.

The project inside the `java-operator-sdk` shows a very simple example where we can see: 
- A controller
- CRD YAML generation 

We can add the `java-operator-sdk` to a Spring Boot project by adding the following dependencies
```
    <dependency>
			<groupId>io.javaoperatorsdk</groupId>
			<artifactId>operator-framework-spring-boot-starter</artifactId>
			<version>3.0.0</version>
		</dependency>
		<dependency>
			<groupId>io.fabric8</groupId>
			<artifactId>crd-generator-apt</artifactId>
			<version>5.12.2</version>
			<scope>provided</scope>
		</dependency>
```

We can create a controller by using annotations and implementing the Reconciler interface: 
```
@ControllerConfiguration
public class ConferenceReconciler implements Reconciler<Conference> {

    @Override
    public UpdateControl<Conference> reconcile(Conference conference, Context context) {
          //Implement logic
        return UpdateControl.updateResource(conference);
    }
}
```

This controller will monitor and reconcile resources of type `Conference`. This means that we need to define this resource and the sub resources like `ConferenceSpec` and `ConferenceStatus`: 

```
@Group("conference.salaboy.com")
@Version("v1")
public class Conference extends CustomResource<ConferenceSpec, ConferenceStatus> implements
        Namespaced {

}

```

The CRD Generator artifact will create the CRD based on these Classes. 

The CRDs can be found in the `target/META-INF/fabric8/` directory after running `mvn package`. 
We need to deal with RBAC configurations for this controller to run inside the cluster, as well as applying the resources to the Cluster. 


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

## Building Controllers with Kubernetes Client Java

You can use the Official Kubernetes Client for Java to build controllers, this will follow the Controller Runtime approach which implement different caching strategies to make sure that our controller is effient and only call the Kubernetes API Server when it needs to:

https://github.com/kubernetes-client/java/blob/master/examples/examples-release-15/src/main/java/io/kubernetes/client/examples/SpringControllerExample.java

This is much low-level but it guarantees that you are following best practices for implementing your controller. 
