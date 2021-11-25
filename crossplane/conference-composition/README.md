# ConferenceInstance Crossplane Composition

In this example, we create a Crossplane Composition to install:
- Conference Platform Services using Helm Provider, as we saw in the [Helm Crossplane example](../helm/README.md)
- PostgreSQL using CloudSQLInstance for GCP as we saw in the [Configuration Package example](../config-pkg/README.md)
- Redis using CloudMemorystoreInstance for GCP as we saw in the [Configuration Package example](../config-pkg/README.md)

![Conference Composition](conference-composition-crossplane.png)

**Note**: You need to have installed the GCP and Helm provider for this composition to work. Check both [Installation](../installing/README.md) and [Helm Example](../helm/README.md) for instructions.

## Creating a conference composition

In the same way that we built our initial abstractions for PostgreSQL and Redis, we can build our conference abstraction using the same 3 files: 
- [crossplane.yaml](crossplane.yaml): This file defines the metadata for our configuration package, as name version and which Crossplane version is required for this package to work. The name for this package `fmtok8s-conference-gcp` makes reference to the fact that it contains a higher level abstraction, allowing us to install a full Conference Platform with GCP resources.
- [composition.yaml](composition.yaml): this file defines how to compose a Conference Platform Instance. It does so, by using the Helm Release introduced in the [Helm Example](../helm/README.md) and the database abstractions that we created with our [Configuration Package example](../config-pkg/README.md).  
- [definition.yaml](definition.yaml): the `definition.yaml` file materialize the simplifications configured by the `composition.yaml` file. In other words, the `definition.yaml` file contains CRD that provides a simple Conference resource that acts as an interface to create Conference Platforms. 

## Building the conference composition package

In order to build and package this configuration package you need to run, inside this directory: 

Build Configuration:

```
kubectl crossplane build configuration
```

Push the OCI package to a registry

```
kubectl crossplane push configuration <USER>/fmtok8s-conference:0.1.0
```

Then install the configuration package into Crossplane: 

```
kubectl crossplane install configuration <USER>/fmtok8s-conference:0.1.0
```

Notice that you can use the composition that I've already built with: 

```
kubectl crossplane install configuration salaboy/fmtok8s-conference:0.1.5
```

Now you can provision conference platforms + their application infrastructure by applying the following yaml file:

```
apiVersion: conferences.fmtok8s.salaboy.com/v1alpha1
kind: ConferenceInstance
metadata:
  name: my-conference
spec:
  parameters: 
    storagePostgresqlGB: 10
    storageRedisGB: 10

```

```
kubectl apply -f resources/conference.yaml
```