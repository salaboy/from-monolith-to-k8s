# Code your CI Pipelines and Run them everywhere with Dagger

In this short tutorial we will be looking at the Dagger Service Pipelines included with each applicaation service. 
These pipelines are implemented in Go using the Dagger Go SDK and take care of building each service, creating a container, publishing it and creating a Helm chart that can be distributed for other teams to use. 

You can find each pipeline at each of the services repository's in a file called `pipeline.go`: 

- [Agenda Service]()
- [Email Service]()
- [Call for Proposals Service]()
- [User Interface]()

To run these pipelines you need to have Go installed and run: 

```
go mod tidy
go run pipeline.go <TASK>
```

If you run `go run pipeline.go full` all the tasks will be executed. Before being able to run all the tasks you will need to make sure that you have all the pre-requisites set, as for pushing to a Container Registry you will need to provide appropriate credentials. 



