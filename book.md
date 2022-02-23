# Continuous Delivery for Kubernetes (Manning)

The contents of this repository had been used to write the [Continuous Delivery for Kubernetes book for Manning](http://mng.bz/jjKP)

This page contains some notes about the chapters and the tutorials included in this repository.

If you have questions, comments or feedback of any kind please create an issue here, or drop me a message in [Twitter @salaboy](http://twitter.com/salaboy). 
If you like the content please give it a Github star!

# Table of Content

- Chapter 1: Cloud-Native Continuous Delivery
- Chapter 2: Cloud-Native applicaiton challenges
- Chapter 3: Service and Environment Pipelines
- Chapter 4: Multi-Cloud Infrastructure
- Chapter 5: Release Strategies
- Chapter 6: Events for Cloud-Native integrations
- Chapter 7: Serverless for Kubernetes
- Chapter 8: TBD
- Chapter 9: TBD 

## Chapter 1: Cloud-Native Continuous Delivery

The first chapter introduce the approach of the book and cover the most important concepts that will be later expanded in the remaining of the book. 
There are two main definitions worth highlighting: What does Cloud-Native means in the context of Kubernetes and Continuous Delivery Goals. 


**Cloud-Native in the context of Kubernetes**: 

A good definition of the term can be found in VMWare site by Joe Beda (Co-Founder, Kubernetes and Principal Engineer, VMware) https://tanzu.vmware.com/cloud-native
“Cloud-Native is structuring teams, culture, and technology to utilise automation and architectures to manage complexity and unlock velocity.”

Where https://12factor.net is a key ingredient of what we understand for Cloud-Native, in Kubernetes we get most of the tools to adhere to these factors, but we need to be aware of them to make sure that we don't design our applications against them. 



**Continuous Delivery Goals**: 
 
“Deliver useful, working software to users as quickly as possible by organising teams to build and deploy in an automated way Cloud-Native applications that run in cloud-agnostic setup.” 

Highly recommended: [Grokking Continuous Delivery - By Christie Wilson - Manning](https://www.manning.com/books/grokking-continuous-delivery)


This chapter also introduce the concept of a walking skeleton, which describes the application linked in the tutorials of this repository. The scenario is a Conference Platform that you can use to organize any kind of conference event with speakers, call for proposals, agenda, notifications, etc.


## Chapter 2: Cloud-Native applicaiton challenges

While looking at the walking skeleton application introduced in chapter 1, chapter 2 focus on describing the challenges that you will face to run this application on top of Kubernetes. 

Starting with defining which Kubernetes flavour you choose and looking into package managers like Helm, this chapter explores the basic Kubernetes resources that you will be using to get the walking skeleton application and up and running. 

This chapter makes reference to the step by step tutorials that can be found [here: https://github.com/salaboy/from-monolith-to-k8s/tree/master/kind](https://github.com/salaboy/from-monolith-to-k8s/tree/master/kind)

Once you get the applicaiton up and running and we review the basics around Kubernetes Deployments, Services, and Ingresses, the chapter goes on describing the most common challenges that developers will face when building this kind of applications.

- **Downtime is not allowed**: If you are building and running a Cloud-Native application on top of Kubernetes and you are still suffering from application downtime, then you are not capitalizing on the advantages of the technology stack that you are using. 
- **Service’s built-in resiliency**: downstream services will go down and you need to make sure that your services are prepared for that. Kubernetes helps with dynamic Service Discovery, but that is not enough for your application to be resilient. 
Dealing with the application state is not trivial: we have to understand each service infrastructural requirements to efficiently allow Kubernetes to scale up and down our services. 
- **Data inconsistent data**: a common problem of working with distributed applications is that data is not stored in a single place and tends to be distributed. The application will need to be ready to deal with cases where different services have different views of the state of the world.
- **Understanding how the application is working (monitoring, tracing and telemetry)**: having a clear understanding on how the application is performing and that it is doing what it is supposed to be doing is essential to quickly find problems when things go wrong. 
- **Application Security and Identity Management**: dealing with users and security is always an after-thought. For distributed applications, having these aspects clearly documented and implemented early on will help you to refine the application requirements by defining “who can do what and when”.  

While following the 12-factor.net principles we will mitigate some of these challenges, we need to consiously design and tackle these challenges to avoid a large number of headaches. 


## Chapter 3: Service and Environment Pipelines

Chapter 3 is about builing Cloud Native application, but going to the next step, of not just building a single service for your applications, but all of them and in an iterative and continuous manner. Automating everything is key for being able to deliver a large amount of services without suffering in the process, hence tools like pipeline engines can help us to automate all the process of building and delivering our services and every new version of them to different environments. 

There are two main concepts covered in this chapter: 
- **Service Pipelines**: the mechanisms to build your services from source until you have everything you need to deploy a runnable version of the service. You will end up having a pipeline definition per service in your application. 
- **Environment Pipelines**: following the GitOps approach, we will have pipelines which their only purpose is to deploy our services to our defined environments. You will end up having one Environment Pipeline per environment that you want to have and these pipelines will take a source of truth as input, such as a git repository which will contain all the information required to set up the environment. When we use this approach, we don't interact with our clusters directly anymore, so no more `kubectl`, we let the Environment Pipeline to promote services and new version from these services to environments.
 
To show how these pipelines can look like I used Tekton, which brings the right constructs to build Cloud-Native pipelines on top of Kubernetes. You can find a tutorial with Tekton which builds on of the services and one environment. 

https://github.com/salaboy/from-monolith-to-k8s/tree/master/tekton



## Chapter 4: Multi-Cloud Infrastructure

When working with Kubernetes, your applications will depend on infrastructure such as Databases, Message Brokers, Load Balancers, API Gateways, Identity Management solutions, etc. All this infrastructure needs to be managed, secured and maintained. Depending on your operation teams, running and managing all this infrastructure is complex and out of the scope for delopment teams. This chapter covers some options to work with infrastructural components inside and outside Kubernetes and how using tools like Crossplane you can manage your Cloud-Provider specific infrastructure by using Kubernetes Resources. 

A tutorial on how to setup the application to use two different databases inside Google Cloud can be found here: https://github.com/salaboy/from-monolith-to-k8s/blob/master/crossplane/config-pkg/README.md

Interesting enough, this example can also be ported to work on any other major Cloud Provider such as Azure AKS or Amazon EKS without changing anything on the application, just Crossplane configuration resources. 


## Chapter 5: Release Strategies

## Chapter 6: Events for Cloud-Native integrations

## Chapter 7: Serverless for Kubernetes

