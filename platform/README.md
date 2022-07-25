# Creating your own self-service platform

With a fast-paced CNCF landspace, it is a full time job to understand, pick and glue projects together to enable your teams with a self-service platform to suit their needs. 

In this document we will be using Kratix, Crossplane, Knative Serving and Knative Functions to demonstrate how a Platform Team can build and curate a set of tools that will enable teams to request using a declarative way. This is just an example of how you can achieve this, and other tools can be used to implement the same behaviours, but some key points that we have tried to cover here are: 
- Self-Service Platform covering two main Personas: Platform Engineers (Platform Team) , Developer Teams (App Team) 
- Developer experience is key to improve productivity, reducing the cognitive load for teams is key to improve efficiency
- Platform Teams can collaborate with the teams to create the right platform for them, using exensible mechanism that can be adapted for more complex needs when needed
- You can achieve all of this by using Open Source projects, but you will need to provide your domain-specific glue


In general, the Platform Team will be in charge of creating the "Platform", developer teams should interact with the "Platform" probably using a portal to create their requests and obtain their environment's credentials. 

![Platform Teams and Developer Teams]
