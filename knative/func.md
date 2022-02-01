# Knative `func` - Serverless on top of Kubernetes

This tutorial shows how to use the `func` CLI in conjunction with Knative and Tekton to build applications using a serverless approach. 
Removing the need from developers to worry about Docker Containers or Kubernetes itself. 

## Pre requisites
- Kubernetes Cluster
- Knative Serving & Eventing installed
- Tekton Installed
- `func` CLI installed


## Use Case
We will be creating a multi-level game, where each level is a function that you can interact with by emitting CloudEvents. 
Functions can be written in any supported language (go, node, python, java, rust, typescript, are we missing your favourite language? get in touch!)


### Creating Level 1

```
func create level-1 -l springboot -t cloudevents
```

