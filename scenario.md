# Scenario (Conference Platform)

Their current application is a Monolith and it looks like this: 
![Monolith Main Site](/imgs/monolith-mainsite.png)
![Monolith Main Backoffice](/imgs/monolith-backoffice.png)

The source code for this application can be [found here](https://github.com/salaboy/fmtok8s-monolith)

The workshop aims to provide the tools, steps, and practices that can facilitate the migration from this application to a Cloud-Native architecture that runs on Kubernetes. In that Journey, we can enable teams to work independently and release their software in small increments while applying some of the principles outlined by the [Accelerate book](https://www.amazon.co.uk/Accelerate-Software-Performing-Technology-Organizations/dp/1942788339/ref=asc_df_1942788339/?tag=googshopuk-21&linkCode=df0&hvadid=311000051962&hvpos=&hvnetw=g&hvrand=13136118265667582563&hvpone=&hvptwo=&hvqmt=&hvdev=c&hvdvcmdl=&hvlocint=&hvlocphy=9072501&hvtargid=pla-446149606248&psc=1&th=1&psc=1). 

![Accelerate](/imgs/accelerate.png)

### Challenges 
In the real world, applications are not that simple. These are some challenges that you might face while doing shift and lift for your Monolith applications:

- **Infrastructure**: if your application has a lot of infrastructure dependencies, such as databases, message brokers, other services, you will need to move them all or find a way to route traffic from your Kubernetes Cluster to this existing infrastructure. If your Kubernetes Cluster is remote, you will introduce latency and security risks which can be mitigated by creating a tunnel (VPN) back to your services. This experience might vary or might be impossible if the latency between the cluster and the services is to high. 

- **More than one process**: your monolith was more than just one application, and that is pushing you to create multiple containers that will have strong dependencies between them. This can be done and most of the time these containers can run inside a Kubernetes Pod if sharing the same context is required.

- **Scaling the application is hard**: if the application hold any kind of state, having multiple replicas becomes complicated and it might require a big refactorings to make it work with multiple replicas of the same running at the same time. 

## Splitting our Monolith into a set of Microservices

Now that we have our Monolith application running in Kubernetes it is time to start splitting it into a set of Microservices. The main reasons to do this are: 
- Enable different teams to work on different parts of this large application
- Enable different services to evolve independently
- Enable different services to be released and deployed independently
- Independently scale services as needed
- Build resiliency into your application, if one service fails not all the application goes down

![Microservices Split](/imgs/microservices-architecture.png)

In order, to achieve all these benefits we need to start simple. The first thing that we will do is add a reverse-proxy which will serve as the main entry point for all our new services. 