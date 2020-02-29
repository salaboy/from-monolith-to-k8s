# From Monolith to K8s

Workshop-style guide for moving from a monolith application to a cloud-native architecture running in Kubernetes.

This guide will take you through an example scenario to move from a Monolith application to a set of Cloud Native microservices running in a Kubernetes Cluster. This workshop highlights the use of certain tools to solve some particular challenges that you will face while going to the cloud. These tools are just suggestions and you should evaluate what fits better to your teams and practices. 

All the projects here are Open Source under the ASL 2.0 License and we welcome Pull Requests and [Issues](http://github.com/salaboy/from-monolith-to-k8s/issues) with more tools additions and suggestions to improve the workshop. 
We encourage people to follow the workshop in their clusters to experience the usage of these tools, their issues and their strengths. 

This workshop is divided into three main sections: 
- Installation and Getting Started
- Running a Cloud Native Conference Application
- Refactoring and improving our applications 

# Installation and Getting Started

This section covers:
 - [Prerequisites]
 - [Tools that we are using during the workshop]
 - [Installing Jenkins X]
 - [Scenario]
 - [CI/CD for our Monolith]

## Prerequisites

- Kubernetes Cluster
  - Tested in GKE Cluster (4 nodes - n1-standard-2)
  - Do you want to test in a different cloud provider and add it to the list? Help is appreciated
- `kubectl` configured. 

## Tools that we are using
- [Kubernetes](http://kubernetes.io)
- [Jenkins X](https://jenkins-x.io)
- [Zeebe](https://helm.zeebe.io)
- Optional (if you want to change code examples and run them locally)
  - [JDK 11+]()
  - [Maven]()

## Installing Jenkins X

First, we will install [Jenkins X](http://jenkins-x.io) in our Kubernetes Cluster. We will use this cluster to build, release and test our Cloud Native applications. 
Jenkins X is divided into a CLI tool `jx` and a set of server-side components. 

We begin by installing the CLI `jx`, if you are running in Mac OSX you can use `brew`

```
brew install jenkins-x/jx/jx
```

Once we have `jx` installed we can run `jx boot` to install the server-side components into the Kubernetes cluster. Remember that `jx` will use the currently configured `kubectl` context, so make sure that you are pointing to the right cluster before running `jx boot`. 

Follow the steps proposed by `jx boot`, for reference these are the options that I've selected in a GKE cluster.  

You can follow a detailed set of instructions from the [Jenkins X Docs page](https://jenkins-x.io/docs/getting-started/). 

> Notice that Jenkins X and this workshop can be executed in any Cloud Provider. 


## Scenario

During this workshop, we will be helping a company that is in charge of providing conference organizers their conference sites and back-office tools for their awesome events. 

Their current application is a Monolith and it looks like this: 
![Monolith Main Site](/imgs/conference-site-main.png)
![Monolith Main Backoffice](/imgs/conference-site-backoffice.png)

The source code for this application can be [found here](https://github.com/salaboy/fmtok8s-monolith)

The workshop aims to provide the tools, steps, and practices that can facilitate the migration from this application to a Cloud-Native architecture that runs on Kubernetes. In that Journey, we can enable teams to work independently and release their software in small increments while applying some of the principles outlined by the [Accelerate book](https://www.amazon.co.uk/Accelerate-Software-Performing-Technology-Organizations/dp/1942788339/ref=asc_df_1942788339/?tag=googshopuk-21&linkCode=df0&hvadid=311000051962&hvpos=&hvnetw=g&hvrand=13136118265667582563&hvpone=&hvptwo=&hvqmt=&hvdev=c&hvdvcmdl=&hvlocint=&hvlocphy=9072501&hvtargid=pla-446149606248&psc=1&th=1&psc=1). 



## CI/CD for our Monolith

When moving to Kubernetes it is quite common to **lift and shift** our monolith applications that are already running outside Kubernertes. This will require to containarize your application and then provide all Kubernetes Manifests to deploy your application into an existing cluster. 

![Monolith](/imgs/monolith-architecture.png)

This exercise pushes us to learn Kubernetes basics concepts such as Deployments, Services, and Ingresses, as well as Docker basics such as how to build and publish a Docker Image and which base Docker Image should we use for our applications. While this is needed for deploying our applications into a running cluster, once we have done these steps for a couple of services/applications, we don't want to do them for 100 services. This is where Jenkins X comes to help us. 

You can find our [monolith application here](http://github.com/salaboy/fmtok8s-monolith). This application is a very basic Spring Boot application which can be started in your local environment (if you have the Java JDK and Maven) installed by running: `mvn spring-boot:run`

If you are running this workshop in your cluster, you can fork the [monolith application here](http://github.com/salaboy/fmtok8s-monolith) repository and then import it to Jenkins X by running:

```
jx import 
```


When we import an application to Jenkins X the following things will happen:
- Our application is decorated with a `Dockerfile` if it doesn't have one, a `jenkins-x.yml` pipeline definition, a `chart` directory containing a [Helm Chart](https://github.com/helm/helm) which contains all the Kubernetes manifests required to deploy our application. 
- Jenkins X will setup webhooks to monitor changes in the application repository. This is why it is required to fork the application code, so Jenkins X has enough rights to set-up webhooks to your repositories. 
- The pipeline defined in the `jenkins-x.yml` file is triggered for the first time in the server-side components. 
- Your application is built, released and promoted to the **Staging Environment**

Some useful commands to track progress are:
- `jx get build logs` -> select your application + ENTER
- `kubectl get pods -n jx-staging` lists all the Pods running in Jenkins X's staging environment

Once the pipeline finishes running you can access your application by running:
`jx get applications` and accessing the URL associated with your application. 

### Challenges 
In the real world, applications are not that simple. These are some challenges that you might face while doing shift and lift for your Monolith applications:
- Infrastructure: if your application has a lot of infrastructure dependencies, such as databases, message brokers, other services, you will need to move them all or find a way to route traffic from your Kubernetes Cluster to this existing infrastructure. If your Kubernetes Cluster is remote, you will introduce latency and security risks which can be mitigated by creating a tunnel (VPN) back to your services. This experience might vary or might be impossible if the latency between the cluster and the services is to high. 
- More than one process: your monolith was more than just one application, and that is pushing you to create multiple containers that will have strong dependencies between them. This can be done and most of the time these containers can run inside a Kubernetes Pod if sharing the same context is required.


# Running a Cloud Native Conference Application

This section covers the following topics: 

- [Splitting our Monolith into a set of Microservices](#splitting-our-monolith-into-a-set-of-microservices)
    - [Introducing an API Gateway](#introducing-an-api-gateway)
    - [Adding a new User Interface](#adding-a-new-user-interface)
    - [Adding more services](#adding-more-services)
      - [Call for Proposals Service](#call-for-proposals-service)
      - [Agenda Service](#agenda-service)
      - [Email Service](#email-service)
    - [Development Flow](#development-flow)
    - [Dealing with infrastructure](#dealing-with-infrastructure)


## Splitting our Monolith into a set of Microservices

Now that we have our Monolith application running in Kubernetes it is time to start splitting it into a set of Microservices. The main reasons to do this are: 
- Enable different teams to work on different parts of this large application
- Enable different services to evolve independently
- Enable different services to be released and deployed independently
- Independently scale services as needed
- Build resiliency into your application, if one service fails not all the application goes down

![Microservices Split](/imgs/microservices-architecture.png)

In order, to achieve all these benefits we need to start simple. The first thing that we will do is add a reverse-proxy which will serve as the main entry point for all our new services. 

## Introducing an API Gateway

If we are going to have a set of services instead of a Monolith application, we will need to deal with routing traffic to each of these new components. In most situations, exposing each of these services outside of our cluster will not be a wise decision. Most of the time, we have a component that is used to aggregate how people access our services from outside the cluster.  

This new component will act as a router between the outside world and our services and you can choose from a set of popular options such as: 
(TBD)

For this workshop, we wanted to use our home-grown component built with [Spring Cloud Gateway](https://spring.io/projects/spring-cloud-gateway), as it gives us the power to tune the routes to our services by coding them in Java or writing these routes in configuration files. 

The source code for our API Gateway can be [found here](http://github.com/salaboy/fmtok8s-api-gateway/)

Once again, you can **fork** and clone this repository and import it to Jenkins X as we did before for the Monolith application.


```
git clone http://github.com/<YOUR USER>/fmtok8s-api-gateway/
cd fmtok8s-api-gateway/
jx import
```

Once again, monitor the pipelines and when the pipeline is finished you should be able to see the new application URL by running:
```
jx get applications
```
and then selecting the pipeline that you want to monitor. 
Wait for the application and environment pipeline to finish to access the application. 

Try to access the API Gateway URL with your browser and see if you can see the new User Interface hosted in this application:
![New User Interface Site](/imgs/conference-microservices-main.png)
![New User Interface Back Office](/imgs/conference-microservices-backoffice.png)

Because we are in the edge, close to our users and outside traffic, the API Gateway serves as the perfect point to host HTML and CSS files that will compose our User Interface.

### Adding a new User Interface

The new user interface will be in charge of consuming REST endpoints which are located now in different services to provide the user the available data. 

We are going to host the new User Interface at the Gateway level as most of these files will reference the API Gateway URL when downloaded to the Client Browser and API Gateways usually provide caching for static files, so the closer these static files are to the user the better. 

The new User Interface look exactly the same as the old one, but in this case to make it more interesting, we will use different colors to the application section to highlight which backend service is in charge of providing data for that section. The User Interface is divided into two main screens, the public **main site** and the **backoffice** which is used by the conference organizers to aprobe/reject proposals and also to send email reminders to people involved in the conference. 

We will also decorate each section with the **version** of the backend service that is serving the requests. 

You can find the logic for the User Interface and the static files inside the [API Gateway source code](https://github.com/salaboy/fmtok8s-api-gateway/blob/master/src/main/java/com/salaboy/conferences/site/DemoApplication.java). 

### Changes from default import
- Memory and CPU allowances: 
- Version Environment Variable: 
- Default Dockerfile with CMD instead use ENTRYPOINT

### Challenges
- Choose the best tool for your team: in this example using the Spring Cloud Gateway made sense as the team already had some Java Knowledge in house, but other reverse proxies provide the same functionality. 
- Running behind a reverse proxy: depending on how flexible your applications, running web applications behind web proxies might require more advanced configurations such as Headers forwarding, tokens forwarding and sometimes path rewrites.
- When we start talking about user interfaces we need to think about authorization and authentication and probably identity management or social logins. This topic is on purpose left out of the workshop as solutions might vary depending on the actual requirements and integrations required by your Cloud Native applications. For OpenID connect with OAuth 2.0 support [Dex is becoming quite popular, you can check it out here](https://github.com/dexidp/dex). 



## Adding more services

In real-life, we start by splitting some peripheral services into microservices to make sure that the core of our application still works. In general, the User Interface can be left untouched when most of the refactorings might happen in the backend. 

Depending on the conference stage, we can start by refactoring out of the Monolith the C4P (Call for Proposals) service, which is in charge of accepting new presentation proposals when the conference is still being organized, while leaving the Agenda untouched still serving users requests. 

It is always recommended to analyze which features and use cases can be used to experiement while thinking about splitting a big Monolith. For the purpose of this workshop we will focus on the **Call for Proposals** flow, hence starting with the Call for Proposals Service. 


## Call for Proposals Service (c4p)

You can find the source code for [this service here](https://github.com/salaboy/fmtok8s-c4p)

> You should **fork** and **jx import** this service as we did with the API Gateway project. 

This service is in charge of handling the logic and the flow for recieving, reviewing and accepting or denying proposals for the conference. Due its responsability it will be in charge of interacting with the Agenda and Email service.
The happy path, or expected flow for this service will be as depictec in the following diagram: 

![C4P Flow](/imgs/c4p-flow.png)

1) Potential Speaker decides to submit a Proposal via the conference site
2) Conference Committee review the proposal and accept or reject it
3) If the proposal is accepted it gets published to the Agenda
4) An email with the result is sent back to the Potential Speaker, informing the decision (approval or rejection)

This service expose the following REST endpoints: 
- GET /info : provide the name and version of the service
- POST / : submit a Proposal
- GET /{id} : get a Proposal by Id
- POST /{id}/decision : decide (approve/reject) a Pending Proposal


### Challenges
- This service is core for the use case that we are trying to cover, as the Call for Papers flow is critical to make sure that we receive, review and make decisions on the proposals sent by potential speakers. This service must work correctly to avoid potential frustration by people submitting valuable proposals. 
- For this example, our services are interacting via REST and our services need to take into considerations that these calls might fail eventually. Retry mechanisms and making sure that our services are idempotent might help to solve these problems. We need to consider also, that if the service interactions fail (due network or services being down) we might end up in an inconsistent state, such as a Potential speaker being approved but not notified. More about this in [Refactoring and improving our applications](#refactoring-and-improving-our-application)


## Agenda Service

The Agenda Service is in charge of hosting the accepted taks for our conference, this service is heavily used by the conference attendees while the conference is on going to check rooms and times for each talk. This service present an interesting usage pattern, as it not going to recieve too many writes (adding new agenda items) as the amount of reads (attendees checking for talks, times and rooms).  

This service expose the following REST endpoints: 
- GET /info : provide the name and version of the service
- POST / : Submit an Agenda Item
- GET / : get all agenda items

### Challenges
- In real-life, a service like this one might justify a separate data store optmized for search and reads. 
- During conference time, we might want to provision more instances (replicas) for this service to serve more traffic
- We might want to consider restricting the POST endpoint when the conference start



## Email Service
The Email Service is an abstraction of a legacy system that you cannot change. Instead of sending emails from the previous services, we encapulated in this case an Email Server behind a REST API. Because we are defining a new API, we can add some domain specific methods to it, such as sending a Conference Notification Email. 

This service expose the following REST endpoints: 
- GET /info : provide the name and version of the service
- POST /notification: send a conference notification email (for Proposals rejections or approvals)
- POST / : send a regular email to an email with Title and Body


## Dealing with infrastructure

When dealing with infrastructural components such as databases, message brokers, identity management, etc, several considerations will influence your decisions:
- Are you running in a Cloud Provider? If so, they probably already offer some managed solutions that you can just use. 
- Do you want to run your own Database instance or Message Broker? If so, look for Helm packages or Kubernetes Operators that help you with managing these components. Scaling a database or a message broker is hard and when you decide to run it yourself (instead of using a managed solution) you add that complexity to managing your applications. 



# Refactoring and improving our applications 

- [Testing and Versioning a Cloud-Native Application](#testing-and-versioning-a-cloud-native-application)
- [Exposing business metrics and insights](#exposing-business-metrics-and-insights)

# Testing and Versioning a Cloud Native Application 

# Exposing business metrics and insights


# References / Links / Next Steps

- [Monolith Application](https://github.com/salaboy/fmtok8s-monolith)
- [API Gateway](https://github.com/salaboy/fmtok8s-api-gateway)