# Conference Controller (MetaController) function

This project was created to demonstrate how to use MetaController and Knative Functions. This controller shows how to monitor a conference application that is running in a separate namespace and trigger some production tests only when the application is healthy. 

This function doesn't require any Kubernetes API Server access, hence it is much easier to code and deploy compared to a normal controller.
The downside, or maybe the advantage, is that we are encouraged more to interact with the data plane instead of the control plane.

For this example, the way to check if the application's services are running is by sending HTTP requests to the endpoints, instead of relying on Kubernetes Resources. 

To create the Deployment that will run the production tests, we use MetaController children definitions. 

To run this project you need: 
- `helm` version 3.8+ installed 
- MetaController installed
  - to install with helm you can run the following commands
  - `HELM_EXPERIMENTAL_OCI=1 helm pull oci://ghcr.io/metacontroller/metacontroller-helm --version=v4.3.7` to fetch the metacontroller chart
  - `kubectl create ns metacontroller` Create a namespace
  - `helm install metacontroller metacontroller-helm-v4.3.7.tgz --namespace metacontroller` install the metacontroller chart
- Knative Serving installed
  - Follow the instructions at https://knative.dev
- `func` CLI installed
  - Follow the instructions at https://github.com/knative-sandbox/kn-plugin-func/blob/main/docs/installing_cli.md

Once you have this setup, there are two main things to do.
- Deploy the function:
  - At the terminal, change to the root of this directory (i.e. `cd kubernetes-controllers/metacontroller/func-conference-controller/`)
  - In the file `func.yaml`, edit the image name to point to your docker registry.   
  - Run `func deploy -v`. This will build, publish and deploy the container image as a Knative Service. Enter your Docker registry credentials at the prompt.
  - Optionally, you can verify that the function has been deployed by running `kubectl get kservice func-conference-controller`.

- Configure the metacontroller to monitor a CRD and then notify our function when a new resource is created. This requires two things:
  - Create a CRD with the type that we want to reconcile, for this example is the `Conference` resource which lives inside the group `metacontroller.conference.salaboy.com` and that we can apply by running `kubectl apply -f config/crd.yaml`
  - Define a MetaController CompositeController where we define that we want to monitor `Conference` resources, we specify which kind of children these resources can have (in this case Deployments) and where (URL) is the function that will do the reconciliation is. You can create these CompositeController by running `kubectl apply -f config/controller.yaml`
  - > Note: If you did not deploy the function to the default namespace, update the url to the function in `config/controller.yaml`.

To test:
- In one terminal window, run `kubectl get pods -w` to watch for pods.
- In a separate terminal window, run `kubectl apply -f config/conference.yaml` to create a resource of type Conference. (Note that you must be in the root of this directory in this window).
- Watch the output in the first window. You should see that a pod is created (the name should be `func-conference-controller-00001-deployment-<UUID>`). Eventually the pod will be terminated. The output may look something like this:
```shell
$ kubectl get pods -w
NAME                              READY   STATUS    RESTARTS   AGE
func-conference-controller-00001-deployment-589ffbc679-q57j8   0/2     Pending   0          0s
func-conference-controller-00001-deployment-589ffbc679-q57j8   0/2     Pending   0          0s
func-conference-controller-00001-deployment-589ffbc679-q57j8   0/2     ContainerCreating   0          0s
func-conference-controller-00001-deployment-589ffbc679-q57j8   1/2     Running             0          2s
func-conference-controller-00001-deployment-589ffbc679-q57j8   2/2     Running             0          5s
func-conference-controller-00001-deployment-589ffbc679-q57j8   2/2     Terminating         0          66s
func-conference-controller-00001-deployment-589ffbc679-q57j8   0/2     Terminating         0          98s
func-conference-controller-00001-deployment-589ffbc679-q57j8   0/2     Terminating         0          98s
func-conference-controller-00001-deployment-589ffbc679-q57j8   0/2     Terminating         0          98s
```

**What happened?**

The metacontroller you created (CompositeController named `metacontroller-conference-controller`) detected the new Conference type resource and sent a request to the function `func-conference-controller`. Knative launched a pod for the function to handle the request and then scaled pod instances back down to zero.




# Generic Function project

Welcome to your new Function project!

This sample project contains a single function based on Spring Cloud Function: `functions.CloudFunctionApplication.uppercase()`, which returns the uppercase of the data passed.

## Local execution

Make sure that `Java 11 SDK` is installed.

To start server locally run `./mvnw spring-boot:run`.
The command starts http server and automatically watches for changes of source code.
If source code changes the change will be propagated to running server. It also opens debugging port `5005`
so a debugger can be attached if needed.

To run tests locally run `./mvnw test`.

## The `func` CLI

It's recommended to set `FUNC_REGISTRY` environment variable.

```shell script
# replace ~/.bashrc by your shell rc file
# replace docker.io/johndoe with your registry
export FUNC_REGISTRY=docker.io/johndoe
echo "export FUNC_REGISTRY=docker.io/johndoe" >> ~/.bashrc
```

### Building

This command builds an OCI image for the function. By default, this will build a JVM image.

```shell script
func build -v                  # build image
```

**Note**: If you want to enable the native build, you need to edit the `func.yaml` file and
set the following BuilderEnv variable:

```yaml
buildEnvs:
  - name: BP_NATIVE_IMAGE
    value: "true"
```

### Running

This command runs the func locally in a container
using the image created above.

```shell script
func run
```

### Deploying

This command will build and deploy the function into cluster.

```shell script
func deploy -v # also triggers build
```

## Function invocation

For the examples below, please be sure to set the `URL` variable to the route of your function.

You get the route by following command.

```shell script
func info
```

Note the value of **Routes:** from the output, set `$URL` to its value.

__TIP__:

If you use `kn` then you can set the url by:

```shell script
# kn service describe <function name> and show route url
export URL=$(kn service describe $(basename $PWD) -ourl)
```

### func

Using `func invoke` command with Path-Based routing:

```shell script
func invoke --target "$URL/uppercase" --data "$(whoami)"
```

If your function class only contains one function, then you can leave out the target path:

```shell script
func invoke --data "$(whoami)"
```

### cURL

```shell script
curl -v "$URL/uppercase" \
  -H "Content-Type:text/plain" \
  -w "\n" \
  -d "$(whoami)"
```

If your function class only contains one function, then you can leave out the target path:

```shell script
curl -v "$URL" \
  -H "Content-Type:text/plain" \
  -w "\n" \
  -d "$(whoami)"
```

### HTTPie

```shell script
echo "$(whoami)" | http -v "$URL/uppercase"
```

If your function class only contains one function, then you can leave out the target path:

```shell script
echo "$(whoami)" | http -v "$URL"
```

## Cleanup

To clean the deployed function run:

```shell
func delete
```
