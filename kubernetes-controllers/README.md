# Extending Kubernetes

This directory contains different Controllers implementations using different tools. 

The main objective for the controllers is to monitor an instance of the Conference application (composed by 4 services), report status back into a Kubernetes resource and triggering production tests, if configured. 

Let's look at the following projects for creating controllers: 
- [Kubebuilder](#kubebuilder)
- [Java Operator SDK](#java-operator-sdk)
- [Kubernetes Client Java](#building-controllers-with-kubernetes-client-java)
- [Metacontroller](#metacontroller)

## Kubebuilder

The [`kubebuilder`](https://book.kubebuilder.io/quick-start.html) directory contains a Kubernetes controller created with `kubebuilder` using Go. 

The instrutions for creating this controller from scratch are: 

- Install the [`kubebuilder` CLI](https://book.kubebuilder.io/quick-start.html)
- `mkdir conference-controller && cd conference-controller`
- `kubebuilder init --domain salaboy.com --repo github.com/salaboy/kubebuilder-conference-controller`
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


## Building Controllers with Kubernetes Client Java

You can use the Official Kubernetes Client for Java to build controllers, this will follow the Controller Runtime approach which implement different caching strategies to make sure that our controller is effient and only call the Kubernetes API Server when it needs to:

https://github.com/kubernetes-client/java/blob/master/examples/examples-release-15/src/main/java/io/kubernetes/client/examples/SpringControllerExample.java

This is much low-level but it guarantees that you are following best practices for implementing your controller. It also remove the man in the middle (fabric8 in this case).

You can find an example controller and the Spring Boot integration here: 
- [Spring Boot Kubernetes Controllers Example](https://github.com/building-k8s-operator/kubernetes-java-operator-sample)

## Metacontroller

With the [MetaController](https://github.com/metacontroller/metacontroller) project you can create Kubernetes controller without the hassle of interacting with the Kubernetes APIs. This is achieved by installing a "MetaController" which allows us to register more focus controllers that doesn't require to interact with the Kubernetes API directly. 

You can define new CRDs and then register a CompositeController to MetaController that will be in charge of notifying your custom controller everytime that one of the resources is created or changed. 

```
apiVersion: metacontroller.k8s.io/v1alpha1
kind: CompositeController
metadata:
  name: metacontroller-conference-controller
spec:
  generateSelector: true
  parentResource:
    apiVersion: metacontroller.conference.salaboy.com/v1
    resource: conferences
  childResources:
  - apiVersion: apps/v1
    resource: deployments
  hooks:
    sync:
      webhook:
        url: http://metacontroller-conference-controller.controllers/
```

This `CompositeController` is monitoring resources with type `conferences.metacontroller.conference.salaboy.com/v1` and it is expected to create `child` resources of type `Deployment`. When a new `Conference` resource is created in the API Server, MetaController will send a POST request to the `webhook` configured above. 

The only requirement for the `Controller` developer is to create a `function` that accept this post request, validate the desired state that comes as a payload and return the current status and if children should be created or not. 

The beauty of this approach is that these functions can run as any Kubernetes Service and they can be written in any programming language. 

![metacontroller diagram](metacontroller-diagram.png)

The fact that MetaController was designed with `functions` in mind makes it really compatible with [Knative Functions](https://github.com/knative-sandbox/kn-func-plugin/).

I have created two separate implementations, one using Spring-Boot and one using a Knative Function also using Spring-Boot. These two implementations are 98% the same, but one can be quickly deployed to a Kubernetes Cluster that have Knative Serving installed in it. 

The main advantages of using Knative Functions are:
- The simplified lifecycle management to create, build and deploy the function to a Kubernetes cluster
- The use of `Spring Cloud Functions` that allows you to define the Function signature using Java's built-in Functional interfaces. 

```
public Function<Map<String, Object>, Mono<Map<String, Object>>> reconcile(){}
```

- Knative Serving will scale down your function to 0 if it is not being used. 

You can find these two projects here: 
- [Spring Boot MetaController Controller](https://github.com/salaboy/from-monolith-to-k8s/tree/main/kubernetes-controllers/metacontroller/conference-controller)
- [Knative Spring Boot Controller Function](https://github.com/salaboy/from-monolith-to-k8s/tree/main/kubernetes-controllers/metacontroller/func-conference-controller)
