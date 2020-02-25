# From Monolith to K8s

Workshop-style guide for moving from a monolith application to a cloud-native architecture running in Kubernetes.

This guide will take you through an example scenario to move from a Monolith application to a set of Cloud Native microservices running in a Kubernetes Cluster. This workshop highlights the use of certain tools to solve some particular challenges that you will face while going to the cloud. These tools are just suggestions and you should evaluate what fits better to your teams and practices. 

All the projects here are Open Source under the ASL 2.0 License and we welcome Pull Requests and [Issues](http://github.com/salaboy/from-monolith-to-k8s/issues) with more tools additions and suggestions to improve the workshop. 
We encourage people to follow the workshop in their own clusters to really experience the usage of these tools.  

# Prerequisites

- Kubernetes Cluster
  - Tested in GKE Cluster (4 nodes - n1-standard-2)
- `kubectl` configured. 


# Tools that we are using
- [Kubernetes](http://kubernetes.io)
- [Jenkins X](https://jenkins-x.io)
- [Zeebe](https://helm.zeebe.io)
- Optional (if you want to change code examples and run them locally)
  - [JDK 11+]()
  - [Maven]()



# Installation phase

First, we will install [Jenkins X](http://jenkins-x.io) in our Kubernetes Cluster. We will use this cluster to build, release and test our Cloud Native applications. 
Jenkins X is divided into a CLI tool `jx` and a set of server-side components. 

We begin by installing the CLI `jx`, if you are running in Mac OSX you can use `brew`

```
brew install jenkins-x/jx/jx
```

Once we have `jx` installed we can run `jx boot` to install the server-side components into the Kubernetes cluster. Remember that `jx` will use the currently configured `kubectl` context, so make sure that you are pointing to the right cluster before running `jx boot`. 

Follow the steps proposed by `jx boot`, for reference these are the options that I've selected in a GKE cluster.  

You can follow a detailed set of instructions from the Jenkins X Docs page. 
> Notice that Jenkins X and this workshop can be executed in any Cloud Provider. 


# Scenario

During this workshop we will be helping a company that is in charge of providing conference organizers their conference sites and back-office tools for their awesome events. 

Their current application is a Monolith and it looks like this: 
![Monolith Main Site](/imgs/conference-site-main.png)
![Monolith Main Backoffice](/imgs/conference-site-backoffice.png)

The source code for this application can be [found here](https://github.com/salaboy/fmtok8s-monolith)

The workshop aims to provide the tools, steps, and practices that can facilitate the migration from this application to a Cloud-Native architecture that runs on Kubernetes. In that Journey, we can enable teams to work independently and release their software in small increments while applying some of the principles outlined by the [Accelerate book](https://www.amazon.co.uk/Accelerate-Software-Performing-Technology-Organizations/dp/1942788339/ref=asc_df_1942788339/?tag=googshopuk-21&linkCode=df0&hvadid=311000051962&hvpos=&hvnetw=g&hvrand=13136118265667582563&hvpone=&hvptwo=&hvqmt=&hvdev=c&hvdvcmdl=&hvlocint=&hvlocphy=9072501&hvtargid=pla-446149606248&psc=1&th=1&psc=1). 

This workshop covers the following steps:
- [CI/CD for our Monolith](#ci-cd-for-our-monolith)
- [Splitting our Monolith into a set of Microservices](#splitting-our-monolith-into-a-set-of-microservices)
    - [Introducing an API Gateway](#introducing-an-api-gateway)
    - [Adding a new User Interface](#adding-a-new-user-interface)
    - [Adding more services](#adding-more-services)
    - [Dealing with infrastructure](#dealing-with-infrastructure)
- [Testing and Versioning a Cloud-Native Application](#testing-and-versioning-a-cloud-native-application)
- [Exposing business metrics and insights](#exposing-business-metrics-and-insights)


# CI/CD for our Monolith

When moving to Kubernetes it is quite common to **lift and shift** our monolith applications that are already running outside Kubernertes. This will require to containarize your application and then provide all Kubernetes Manifests to deploy your application into an existing cluster. 

![Monolith](/imgs/monolith-architecture.png)

This exercise pushes us to learn Kubernetes basics concepts such as Deployments, Services, and Ingresses, as well as Docker basics such as how to build and publish a Docker Image and which base Docker Image should we use for our applications. While this is needed for deploying our applications into a running cluster, once we have done these steps for a couple of services/applications, we don't want to do them for 100 services. This is where Jenkins X comes to help us. 

You can find our [monolith application here](http://github.com/salaboy/fmtok8s-monolith). This application is a very basic Spring Boot application which can be started in your local environment (if you have the Java JDK and Maven) installed by running: `mvn spring-boot:run`

If you are running this workshop in your cluster, you can fork the [monolith application here](http://github.com/salaboy/fmtok8s-monolith) repository and then import it to Jenkins X by running:

```
jx import 
```

or (if you haven't cloned your fork)

```
jx import --url http://... 
```

When we import an application to Jenkins X the following things will happen:
- Our application is decorated with a `Dockerfile` if it doesn't have one, a `jenkins-x.yml` pipeline definition, a `chart` directory containing a [Helm Chart](https://github.com/helm/helm) which contains all the Kubernetes manifests required to deploy our application. 
- Jenkins X will setup webhooks to monitor changes in the application repository. This is why it is required to fork the application code, so Jenkins X have enough rights to setup webhooks to your repositories. 
- The pipeline defined in the `jenkins-x.yml` file is triggered for the first time in the server-side components. 
- Your application is built, released and promoted to the **Staging Environment**

Some useful commands to track progress are:
- `jx get build logs` -> select your application + ENTER
- `kubectl get pods -n jx-staging` lists all the Pods running in Jenkins X's staging environment
- 

Once the pipeline finishes running you can access your application by running:
`jx get applications` and accessing the URL associated with your application. 

### Challenges 
In the real world, applications are not that simple. These are some challenges that you might face while doing shift and lift for your Monolith applications:
- Infrastructure: if you application has a lot of infrastructure dependencies, such as databases, message brokers, other services, you will need to move them all or find a way to route traffic from your Kubernetes Cluster to this existing infrastructure. If your Kubernetes Cluster is remote, you will introduce latency and security risks which can be mitigated by creating a tunnel (VPN) back to your services. This experience might vary or might be impossible if the latency between the cluster and the services is to high. 
- More than one process: your monolith was more than just one application, and that is pushing you to create multiple containers that will have strong dependencies between them. This can be done and most of the time these containers can run inside a Kubernetes Pod if sharing the same context is required.
- (TBD)

# Splitting our Monolith into a set of Microservices

Now that we have our Monolith application running in Kubernetes it is time to start splitting it into a set of Microservices. The main reasons to do this are: 
- Enable different teams to work on different parts of this large application
- Enable different services to evolve independently
- Enable different services to be released and deployed independently
- Independently scale services as needed
- Build resiliency into your application, if one service fails not all the application goes down

![Microservices Split](/imgs/microservices-architecture.png)

In order, to achieve all these benefits we need to start simple. The first thing that we will do is add a reverse-proxy which will serve as the main entry point for all our new services, including the monolith. 

## Introducing an API Gateway

If we are going to have a set of services instead of a Monolith application, we will need to deal with routing traffic to each of these new components. In most situations, exposing each of these services outside of our cluster will not be a wise decision. Most of the time, we have a component that is used to aggregate how people access our services from outside the cluster.  

This new component will act as a router between the outside world and our services and you can choose from a set of popular options such as: 
(TBD)

For this workshop, we wanted to use our home-grown component built with [Spring Cloud Gateway](https://spring.io/projects/spring-cloud-gateway), as it gives us the power to tune the routes to our services by coding them in Java or writing these routes in configuration files. 

The source code for our API Gateway can be [found here](http://github.com/salaboy/fmtok8s-api-gateway/)

Once again, you can fork and clone this repository and import it to Jenkins X as we did before for the Monolith application.

```
jx import
```

Once again, monitor the pipelines and when the pipeline is finished you should be able to see the new application URL by running:
```
jx get applications
```

Try to access the API Gateway URL with your browser and see if you can see the new User Interface hosted in this application:
![New User Interface Site](/imgs/conference-microservices-main.png)
![New User Interface Back Office](/imgs/conference-microservices-backoffice.png)


### Adding a new User Interface

The new user interface will be in charge of consuming REST endpoints which are located now in different services to provide the user the available data. 

We are going to host the new User Interface at the Gateway level as most of these files will reference the API Gateway URL when downloaded to the Client Browser and API Gateways usually provide caching for static files, so the closer these static files are to the user the better. 

The new User Interface can look exactly the same as the old one, but in this case to make it more interesting, we will use different colors to the application section to highlight which backend service is in charge of providing data for that section. 

We will also decorate each section with the **version** of the backend service that is serving the requests. 

### Changes from default import
- Memory and CPU allowances: 
- Version Environment Variable: 
- (Java) JDK version

### Challenges
- Choose the best tool for your team: in this example using the Spring Cloud Gateway made sense as the team already had some Java Knowledge in house, but there are other reverse proxies that can provide the same functionality. 
- Runing behind a reverse proxy: depending on how flexible your applications, running web applications behind web proxies might require more advanced configurations such as Headers forwarding, tokens forwarding and sometimes path rewrites. 



## Adding more services

In the real life, we start by splitting some periferial services into microservices to make sure that the core of our application still works. 
Depending on the conference stage, we can start by refactoring out of the Monolith the C4P (Call for Proposals) service, which is in charge of accepting new presentation proposals when the conference is still being organized, while leaving the Agenda un touched and still serving users requests. 

### Call for Proposals Service (c4p)

You can find the source code for [this service here]()

This service is in charge of handling the logic and the flow for recieving, reviewing and accepting or denying proposals for the conference. 

### Agenda Service
(TBD)

### Email Service
(TBD)


## Dealing with infrastructure

When dealing with infrasctural components such as databases, message brokers, identity management, etc, there are several considerations that will influence your decisions:
- Are you running in a Cloud Provider? If so, they probably already offer some managed solutions that you can just use. 
- Do you want to run your own Database instance or Message Broker? If so, look for Helm packages or Kubernetes Operators that help you with managing these components. Scaling a database or a message broker is hard and when you decide to run it yourself (instead of using a managed solution) you add that complexity to managing your applications. 



# Testing and Versioning a Cloud Native Application 

# Exposing business metrics and insights


# References / Links / Next Steps

- [Monolith Application](https://github.com/salaboy/fmtok8s-monolith)
- [API Gateway](https://github.com/salaboy/fmtok8s-api-gateway)