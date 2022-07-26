# Creating your own self-service multi-cluster platform

With a fast-paced [CNCF landspace](https://landscape.cncf.io/), it is a full time job to understand, pick and glue projects together to enable your teams with a self-service platform to suit their needs. 

In this tutorial we will be using Crossplane, Tekton, Knative Serving, ArgoCD and Knative Functions to demonstrate how a Platform Team can build and curate a set of tools that will enable teams to request using a declarative way. This is just an example of how you can achieve this, and other tools can be used to implement the same behaviours, but some key points that we have tried to cover here are: 
- Self-Service Platform covering two main Personas: Platform Engineers (Platform Team) , Developer Teams (App Team) 
- Developer experience is key to improve productivity, reducing the cognitive load for teams is key to improve efficiency
- Platform Teams can collaborate with the teams to create the right platform for them, using exensible mechanism that can be adapted for more complex needs when needed
- You can achieve all of this by using Open Source projects, but you will need to provide your domain-specific glue.  


Most of what is covered here revolves around: 
- Which tools do the Platform Teams use to create a self-service, API-driven platform that developer can use
- How do developers consume these APIs in a self-service manner and what practices are promoted by the platform



![Platform Teams and Developer Teams]


## One Cluster is not enough

If your organization is large enough to have more than one development team building services you face the need to isolate their environments so they don't block each other. If these teams are [extending Kubernetes](https://github.com/salaboy/from-monolith-to-k8s/tree/main/kubernetes-controllers) and creating their own Custom Resource Definitions or defining cluster-wide configurations it's wise to isolate their work in separete clusters to avoid any issues.

Using Kubernetes `namespaces` for isolation is usually not enough, but fully creating a new cluster per team is not cost-efficient for all the use cases. 

If we are building a Platform we should understand these requirements and encode them in a way that development teams can request the creating of new environments for them to work and not worry about where other teams are working. 

If we want to contemplate the entire spectrum from `namepsace` isolation to `cluster` isolation one thing remains consistent, we need to automate how these environment gets created in a way that is abstracted from the development teams and that is offered in a self-service approach. Topics like security configurations, identity management, cluster provisioning and installed tools required for their development tasks shouldn't be delegated to developers who should be focused on writing features. 

@TODO: create a table to compare namepsace and cluster isolation

We also want to encode these decisions at the platform level, so we need they right tools that allows us (Platform builders) to encode these decisions in a way that we can include more projects of change these definitions if our context change. 


In this section we will take a look at two tools which can help you to provide some automation around these topics: 
- Crossplane
- VCluster

With Crossplane we can create a Kubernetes Cluster and other cloud resource in any Cloud Provider by creating Compositions. These compositions allows us to create higher-level constructs to expose only the parameters relevant to our development teams. They don't even need to know that they are creating a Kubernetes Cluster, they only need to know where to connect to deploy their applications. These compositions can also take care of security configurations and creating other cloud resources such as databases or message brokers. 

With VCluster we can create new Virtual Clusters inside a single Kubernetes Cluster. This gives us flexiblity to have isolated API Servers for each Virtual Cluster but always inside the same `host` Kubernetes Cluster. This reduces the management around creating and provisioning clusters, but more importantly, we avoid creating new dedicated Control Planes that we will need to pay for and are not running services that 


## APIs and Self-service

As Platform builders, we need to craft a contract between the Platform and its users (Dev teams) and there are several approaches to do this. In this section we will evaluate two different approaches that are becoming quite popular in the Cloud-Native space: 
- Crossplane Compositions as a way to rely on the Kubernetes APIs
- Kratix as the API builder
- Compositions as Code with Pulumi

These 3 approaches try to tackle the same topic give us flexibility at different levels, so it is worth looking at their advantages to contrast these solutions and maybe used them as complimentary technologies. 


- Crossplane compositions as they are today doesn't allow conditionals which might be required for complex scenarios. The big advantage with Crossplane is that the language to define these compositions is just Kubernetes resources, hence we can use all the tools available in the Kubernetes ecosystem to manage these definitions. These compositions are also automatically discoverable by users that can interact with the Kubernetes APIs as they are just Kubernetes resources. 
- Pulumi allows you to encode these compositions using your favourite programming language. If you are using Java or Go the APIs and logic used to build these compositions can be defined using the tools that your development teams are already using. At the end of the day, a composition can be packaged as just a library that you can manage as any other source code artifact. The drawback of this approach is that it doesn't use the Kubernetes APIs to describe these compositions, hence Kubernetes is not aware of what is available. 
- Kratix only allows us to define an API by creating new CRDs and Controllers for these new resources, what happens behind these APIs is totally up to you. While this approach is clean and removes the need for creating your custom controllers, as a platform team you will need to choose which tools are you going to use to implement the behaviours of these new resources (abstractions). 



## Domain-specific Opinions

- One or multiple clusters, which tools? 
- Pipelines and supply chain
- Release strategies
- Developer Experience



## TODO:
- Demo: Crossplane composition as an abstraction for our Conference Application
- Pulumi APIs to compose configuraitons
- Kratix API for Supply Chain 
- Don't forget about: Developer Experience


