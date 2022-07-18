# Conference Controller using Java Operator SDK

You can find more information about this controller in the following repository: [http://github.com/salaboy/from-monolith-to-k8s](http://github.com/salaboy/from-monolith-to-k8s)

This controller is built using the Java Operator SDK inside a Spring Boot application hence most of the lifecycle of the application can be handled using Spring Boot.

To run this controller locally you need to:
- `mvn package -DskipTests` to build the artifact and the kubernetes resources. You need to skip the tests so the artifacts that are required by the controller are generated. 
- `kubectl apply -f target/classes/META-INF/fabric8`, if one of the resources applied failed, that's ok
- `mvn spring-boot:run` to start the controller locally
- `kubectl create ns java-operator-sdk` to have the resource isolated from other resources.
- `kubectl apply -f config/conference.yaml -n java-operator-sdk`

When running the controller locally, it will fail to monitor the application that is running in a remote cluster, hence you might want to run the controller inside a cluster. 

## Running the controller inside a Kubernetes Cluster

Running controllers is a complex topic, mostly because we need to create a container, publish it in a remote registry that the Kubernetes Cluster can access to fetch the container image
and then make sure that we have the correct RBAC configurations for the controller to access the resources that it needs from inside the cluster. 

You can deploy the controller into a remote Kubernetes cluster by running the following steps:
- `mvn spring-boot:build-image` to create the container
- Tag and push created container (replace `salaboy` for your registry user):
```bash
docker tag java-operator-sdk-conference-controller:0.0.1-SNAPSHOT salaboy/java-operator-sdk-conference-controller:java-operator-sdk
docker push salaboy/java-operator-sdk-conference-controller:java-operator-sdk
```
- `kubectl apply -f config/controller.yaml -n java-operator-sdk` to deploy the controller to the `java-operator-sdk` namespace. 

Then you can create conferences with:
- `kubectl apply -f config/conference.yaml -n java-operator-sdk`

If you have a conference running in the specified namespace (inside the `conference.yaml` resource) the service checks will work and a new deployment should be created inside the `java-operator-sdk` namespace.
