# Simple Go program to compare with Java Frameworks

This very very very simple project aims to show the syntax and main constructs that we can use to build a simple
webserver.

## Initializing the project

To create this project I've created a new directory, and then I've run:

```bash
go mod init github.com/salaboy/from-monolith-to-k8s/kubernetes-controllers/javavsgo/go/conference-service
```

This project is using [Go Modules](https://go.dev/ref/mod) which is the dependency management tool which is built-in in
Go.

## Struct

In Java this will relate to what we know as Plain Old Java Objects (POJO) or nowadays to a Java Record

```go
type Conference struct {
Id    string    `json:"id"`
Name  string    `json:"name"`
Where string    `json:"where"`
When  time.Time `json:"when"`
}
```

There is no need for getters of setters, and we control the visibility of the fields using caps for each field first
letter. All previous fields are public. The struct itself is also public because it was defined with capital `C`.

## Structs and Behaviour

If we want to add functionality to our `struct`s we can make a function part of the struct, but not inside the struct
definition.

Here a simple example of how you attach a `func`tion to a `struct`:

```go
type ConferenceStore struct {
}

func (*ConferenceStore) read() []Conference {
    return []Conference{
        {
            Id:    "123",
            Name:  "JBCNConf",
            Where: "Barcelona, Spain",
            When:  time.Date(2022, time.July, 18, 0, 0, 0, 0, time.UTC),
        },
        {
            Id:    "456",
            Name:  "KubeCon",
            Where: "Detroit, USA",
            When:  time.Date(2022, time.October, 24, 0, 0, 0, 0, time.UTC),
        },
    }
}
```

To use the `read()` function you need an instance of the `ConferenceStore` struct:

```go
store := ConferenceStore{}
conferences := store.read()
```

## Unit Testing

Unit Tests in `Go` go right besides the code that they are supposed to test. `go test` will look for all the files
called `*_test.go` and treat them as tests.

Here is a simple test to test our `ConferenceBuilder.build()` function:

```go
func TestConferenceBuilder(t *testing.T){
builder := ConferenceBuilder{}
got := builder.build()
want := Conference{
Id:    "123",
Name:  "JBCNConf",
Where: "Barcelona, Spain",
When:  time.Date(2022, time.July, 18, 0, 0, 0, 0, time.UTC),
}

if got != want {
t.Errorf("got %q, wanted %q", got, want)
}
}
```

## Creating a simple REST endpoint for our Conference Builder

You can build a simple web server and expose the endpoint using basic Go constructs or you can use a library. For this
example and because we want to compare the Go Ecosystem we will be using [`Gorilla MUX`](https://github.com/gorilla/mux)
to create a simple server and expose our simple GET endpoint.

To install a new library into your project you will end up doing something like this:

```bash
go get -u github.com/gorilla/mux
```

This tells Go Modules to add the dependency to the `go.mod` file. You can see now in the `go.mod` file a line like this:

```
require github.com/gorilla/mux v1.8.0
```

Basically, the `mux` library latest version is `v1.8.0`

## Creating a WebServer and a new REST endpoint with MUX

With MUX we need to create a new `Router` to define which endpoints and which `Handler`s will be in charge of the
requests:

```go
r := mux.NewRouter()
r.HandleFunc("/conferences", ConferencesHandler ).Methods(http.MethodGet)
http.Handle("/", r)
log.Fatal(http.ListenAndServe(":8080", nil))
```

The `http.ListenAndServe` comes built in with Go.

Take a look into the `main_test.go` file to check the test for the endpoint.

## Kubernetes and `ko`

First things first, if we want to run our simple service in Kubernetes we will need to add a `health` endpoint.

The Gorilla Mux doc provides us with a very simple `health` endpoint:

```
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, `{"alive": true}`)
}
```

Check the `main_test.go` for the unit test for this Handler.

Next, we need to build a container and have some YAML files so we can deploy our service to our Kubernetes Cluster. That
is where [`google/ko`](https://github.com/google/ko) comes really handy.

Once you install `ko` you might want to setup the docker registry that you want to use by setting the following Env Var:

```
export KO_DOCKER_REPO=salaboy
```

This will use Docker Hub, but check the `ko` website to configure other registries such as the GitHub Container
Registry (ghcr) or Google Cloud Registry (gcr).

To build our application container we can run now:

```
ko build main.go
```

This will build and push the container image to our configured registry. For an application like this one you will see
that the image name is basically `main.go`

```
salaboy/main.go-7ddfb3e035b42cd70649cc33393fe32c@sha256:870201bece24bf8f841100be55575975219d04c374f6f7216365c9f2f9c901c4
```

Now that we have a container we can do a couple of things. First we can run this container using docker, but `ko` allows
us to make sure that what we run is a container that reflect the code that we are writing. The following command will
build, push and then run the latest version of our container.

```
docker run -p 8080:8080 $(ko build main.go)
```

Do you want a multi-platform container:

```
ko build main.go --platform=linux/amd64,linux/arm64
```

Now let's deploy our service to Kubernetes, `ko apply` to the rescue:
`ko apply` will run `ko build` which will build and push the container image to the configured registry and then
use `ko resolve` to replace the container image name (and SHA) inside the Kubernetes YAML files inside the `config/`
directory.
`ko` uses the prefix `ko://` to find and replace the container image and uses the `main.go` to know which image apply.

Notice that I've manually created these YAML files.

```
ko apply -f config/
```

After running this command, we can port-forward to the service and interact with it inside our Kubernetes Cluster:

```
kubectl port-forward svc/conference-service 8080:80
```

If we make a change in the service, we can just run again:

```
ko apply -f config/
```

A new pod will be created, as there is a new SHA for the container image. 





