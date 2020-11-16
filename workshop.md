# Workshop 

During this workshop you will deploy a Cloud-Native application, inspect it, change its configuration to use different services and 
play around with it to get familiar with Kubernetes and Cloud-Native tools that can help you to be more effective in your cloud journey. 

During this workshop you will be using GKE (Managed Kubernetes Engine inside Google Cloud) to deploy a complex application composed by multiple services. But none of the applications or tools used are tied in any way to Google infrastructure, meaning that you can run these steps in any other Kubernetes provider, as well as in an On-Prem Kubernetes installation. 

The main goal of the workshop is to guide you step by step to work with an application that you don't know but that will run on a real infrastructure (in contrast to run software in your own laptops). Due the time constraints, the workshop is focused on getting things up and running, but it opens the door for a lot of extensions and experimentation, that we encourage. 

# Creating accounts and Clusters
During this workshop you will be using a **Kubernetes Cluster** and a **Camunda Cloud Zeebe Cluster** for Microservices orchestration. You need to setup these accounts and create these clusters early on, so they are ready for you to work for the reminder of the workshop. 

**Notice: you can click in the images to expand them**

1) [Login here to Google Cloud](http://console.cloud.google.com) (if you are a QCon Plus attendee, we will provide you with one)
<details>
  <summary>Creating a Kubernetes Cluster (Click to Expand)</summary>

We recommend to use an **Incognito Window** in **Chrome** (File -> New Incognito Window) to run the following steps, as with that you will avoid having issues with your personal **Google Account**

Once you are logged in, you will be asked to accept the terms and continue: 
<img src="workshop-imgs/00-accept-terms-and-continue.png" alt="Terms" width="500px">

Once the terms are accepted, it is **extremely important** that you select the correct project to work on. On the top bar, there is a project dropdown that opens the project list. You need to click into the **QCon SF 2020 ...** project to select it.

<img src="workshop-imgs/01-select-qcon-project.png" alt="Select Project" width="500px">

Once the project is selected, you can create new **Kubernetes Clusters* by switching to the **Kubernetes Engine** section from the left hand side menu:

<img src="workshop-imgs/02-go-to-kube-engine.png" alt="Kubernetes Engine section" height="400px">

Once in the **Kubernetes Engine** section you will notice that Google Cloud will initialize the **Kubernetes APIs** for us, you need to wait for this to finish:

<img src="workshop-imgs/03-gcp-enabling-kube-apis.png" alt="GCP Enabling Kube APIS" width="500px">

Once the Kubernetes APIs are enabled, you will be able to create a new Kubernetes Cluster by hitting the create button, that now should be enabled:

<img src="workshop-imgs/04-gcp-create-kube-cluster.png" alt="Create Cluster" width="500px">

You will be creating a **3 Nodes (n2-standard-4)** Cluster. The first step is to name your cluster, use the name **workshop**

<img src="workshop-imgs/05-cluster-basics-name-workshop.png" alt="Name it workshop" width="500px">

In the **Node Pools -> Default pool** section (on the left hand side menu) check that the **Kubernetes Master** version is **1.16+** (which should be the default) and that the number of nodes is **3**:

<img src="workshop-imgs/06-cluster-pool.png" alt="Name it workshop" width="500px">

Finally, you need to define which kind of computers will be provisioned for your cluster, for doing this switch to the **Node** section in the left hand side menu and select **N2** in the **Series** dropdown and **n2-standard-4** in the **Machine Type** dropdown:

<img src="workshop-imgs/07-cluster-nodes.png" alt="Name it workshop" width="500px">

Finally, hit **Create** at the bottom of the screen. This triggers the provisioning of the machines required for your cluster and the setup process required by Kubernetes. 

This creation process takes several minutes, you will see the loading icon right beside your cluster name:
<img src="workshop-imgs/08-wait-for-cluster.png" alt="Waiting" width="500px">

This will take some minutes, so you can move forward to **Camunda Cloud Account and Cluster** while the Kubernetes Cluster is being created.
  
</details>  
  
2) [Create a Camunda Cloud Account and Cluster](https://accounts.cloud.camunda.io/signup?campaign=workshop) 
<details>
  <summary>Login into your account and create a Cluster (Click to Expand)</summary>

**Fill the form** to create a new account, you will need to use your email to confirm your account creation. You will be using Camunda Cloud for **Microservices Orchestration** ;)  
<img src="workshop-imgs/13-create-camunda-cloud-account.png" alt="Create Camunda Cloud Account" width="500px">

Check your inbox to **Activate your account** and follow the links to login, after confirmation:

<img src="workshop-imgs/14-activate-your-account.png" alt="Activate" width="500px">

**Once activated, Login with your credentials** and let's create a new **Zeebe Cluster**, you will be using this cluster later on in the workshop, but it is better to set it up early on. 

<img src="workshop-imgs/15-create-a-new-zeebe-cluster.png" alt="Create Cluster" width="500px">

**Create a new cluster** called `my-cluster`:

<img src="workshop-imgs/16-call-it-my-cluster.png" alt="My Cluster" width="500px">

Disregard, creating a model if you are asked to and just close the popup:

<img src="workshop-imgs/17-disregard-creating-model.png" alt="Close popup" width="500px">

**Your cluster is now being created:**

<img src="workshop-imgs/18-cluster-is-being-created.png" alt="Cluster is being created" width="500px">


</details>  

Let's switch back to Google Cloud to setup your Kubernetes Cluster to start deploying our Cloud-Native Applications!

# Connecting to your Kubernetes Cluster and installing Knative

During this workshop, you will be using Cloud Shell to interact with your Kubernetes Cluster, this avoids you setting up tools in your local environment and it provides quick access to the cluster resources. 


Once the Kubernetes cluster is created and "green", you will connect and iteract with it using **Cloud Shell**, a terminal that runs inside a Debian machine which comes with pre-installed tools like: `kubectl` and `helm`. 

Click the **Connect** button 

<img src="workshop-imgs/09-cluster-green.png" alt="Cluster Green Connect" width="500px">

Then find the **Run in Cloud Shell** button, which will provision a new instance of Cloud Shell for you to use:

<img src="workshop-imgs/10-connect-to-cluster-with-cloud-shell.png" alt="Cloud Shell" width="500px">

Once **Cloud Shell** is provisioned, notice that you will need to hit **Enter** to actually connect with the **workshop** cluster:

<img src="workshop-imgs/11-connect-from-cloud-shell.png" alt="Cloud Shell" width="700px">

Because you will be using the `kubectl` and `helm` commands a lot during the next couple of hours we recommend you to create the following aliases:

``` bash
alias k=kubectl
alias h=helm
```
Now instead of typing `kubectl` or `helm` you can just type `k` and `h` respectivily. 

You can now type inside Cloud Shell: 
``` bash
k get nodes
```

You should see something like this: 

<img src="workshop-imgs/12-tests-kubectl.png" alt="Cloud Shell" width="700px">


Next step you will install Knative Serving and Knative Eventing. The Cloud-Native applications that you will deploy in later steps were built having Knative in mind. 

- Knative Serving

If you have the previous aliases set up, you can copy the entire block and paste it Cloud Shell

``` bash
k apply --filename https://github.com/knative/serving/releases/download/v0.18.0/serving-crds.yaml
k apply --filename https://github.com/knative/serving/releases/download/v0.18.0/serving-core.yaml
k apply --filename https://github.com/knative/net-kourier/releases/download/v0.18.0/kourier.yaml
k patch configmap/config-network \
  --namespace knative-serving \
  --type merge \
  --patch '{"data":{"ingress.class":"kourier.ingress.networking.knative.dev"}}'
k apply --filename https://github.com/knative/serving/releases/download/v0.18.0/serving-default-domain.yaml

```

You can check that the installation worked out correctly by checking that all the Knative Serving pods are running:
```bash
k get pods -n knative-serving
```

You should see something like this:

<img src="workshop-imgs/25-knative-serving-test.png" alt="KNative Serving Test" width="700px">

- Knative Eventing

``` bash
k apply --filename https://github.com/knative/eventing/releases/download/v0.18.0/eventing-crds.yaml
k apply --filename https://github.com/knative/eventing/releases/download/v0.18.0/eventing-core.yaml
k apply --filename https://github.com/knative/eventing/releases/download/v0.18.0/in-memory-channel.yaml
k apply --filename https://github.com/knative/eventing/releases/download/v0.18.0/mt-channel-broker.yaml
k create -f - <<EOF
apiVersion: eventing.knative.dev/v1
kind: Broker
metadata:
 name: default
 namespace: default
EOF

```

You can check that the installation worked out correctly by checking that all the Knative Eventing pods are running:
``` bash
k get pods -n knative-eventing
```

You should see something like this: 

<img src="workshop-imgs/26-knative-eventing-test.png" alt="KNative Eventing Test" width="700px">

Now, you have everything ready to deploy your Cloud-Native applications to Kubernetes. 

# Version 1: Deploying a Cloud-Native Application

In this section you will be deploying a Conference Cloud-Native application composed by 4 simple services. 

With Knative installed you can proceed to install the first version of the application. You will do this by using [**Helm**](http://helm.sh) a Kuberenetes Package Manager. As with every package manager you need to add a new `Helm Repository` where the **Helm packages/charts** for the workshop are stored. 

You can do this by runnig the following commands: 

``` bash
h repo add workshop http://chartmuseum-jx.35.222.17.41.nip.io
h repo update

```

Now you are ready to install the application by just running the following command:
``` bash
h install fmtok8s workshop/fmtok8s-app

```
You should see something like this (ignore the warnings):

<img src="workshop-imgs/27-helm-repo-add-update-install-v1.png" alt="Helm install" width="700px">

The application [Helm Chart source code can be found here](https://github.com/salaboy/fmtok8s-app/).

You can check that the application running with the following two commands:

- Check the pods of the running services with: 
``` bash
k get pods
```

- You can also check the Knative Services with: 
```
k get ksvc
```

You should see that pods are being created or they are running and that the Knative Services were created, ready and have an URL:

<img src="workshop-imgs/28-k-getpods-kgetksvc.png" alt="kubectl get pods and ksvcs" width="700px">

As soon all the pods are running and the services are ready you can copy and paste the `fmtok8s-api-gateway` URL into a different tab in your browser to access the application `http://fmtok8s-api-gateway.default.XXX.xip.io`

<img src="workshop-imgs/agenda-screen.png" alt="Conference Agenda" width="500px">

Now you can go ahead and:
1) Submit a proposal by clicking the Submit Proposal button in the main page

2) Go to the back office (top right link) and Approve or Reject the proposal

<img src="workshop-imgs/backoffice-screen.png" alt="Conference BackOffice" width="500px">

3) Check the email service to see the notification email sent to the potential speaker, this can be done with 
``` bash
k get pods
```
Where you should see the Email Service pod:

<img src="workshop-imgs/59-v1-get-pods-email-highlighted.png" alt="Conference BackOffice" width="500px">

And then you can tail the logs by running:
``` bash
k logs -f fmtok8s-email-<YOUR POD ID> user-container
```
You should see the service logs being tailed, you can exit/stop taling the logs with `CTRL+C`.

<img src="workshop-imgs/60-email-service-spring-boot-started.png" alt="Conference BackOffice" width="500px">

And if you **approved** the submitted proposal you should also see something like this: 

<img src="workshop-imgs/61-email-service-tail-logs-approved.png" alt="Conference BackOffice" width="500px">

4) If you approved the proposal, the proposal should pop up in the Agenda (main page) of the conference. 

<img src="workshop-imgs/62-proposal-in-agenda.png" alt="Conference BackOffice" width="500px">

Let's take a deeper look on what you just did in this section. 

# Understanding your application

In the previous section you installed an application using `Helm`. 

For this example, there is a parent **Helm Chart** that contains the configuration for each of the services that compose the application. 
You can find each service that is going to be deployed inside the `requirements.yaml` file defined [inside the chart here](https://github.com/salaboy/fmtok8s-app/blob/master/charts/fmtok8s-app/requirements.yaml).

This can be extended to add more components if needed, like for example adding application infrastructure components such as Databases, Message Brokers, ElasticSearch, etc. (Example: [ElasticSearch](https://github.com/elastic/helm-charts), [MongoDB](https://artifacthub.io/packages/helm/bitnami/mongodb) and [MySQL](https://artifacthub.io/packages/helm/bitnami/mysql), [Kafka](https://artifacthub.io/packages/helm/bitnami/kafka) charts). 

The configuration for all these services can be found in the [`value.yaml` file here](https://github.com/salaboy/fmtok8s-app/blob/master/charts/fmtok8s-app/values.yaml). This `values.yaml` file can be overriden as well as any of the settings from each specific service when installing the chart, allowing the chart to be flexible enough to be installed with different setups. 

There are a couple of configurations to highlight for this version which are:
- [Knative Deployments are enabled](https://github.com/salaboy/fmtok8s-app/blob/master/charts/fmtok8s-app/values.yaml#L6), each service Helm Chart enable us to define if we want to use a Knative Service or a Deployment + Service + Ingress type of deployment. Because we have Knative installed, and you want to leverage Knative 
- Both the [`C4P` service](https://github.com/salaboy/fmtok8s-app/blob/master/charts/fmtok8s-app/values.yaml#L16) and the [`API Gateway` service](https://github.com/salaboy/fmtok8s-app/blob/master/charts/fmtok8s-app/values.yaml#L7) need to know where the other services are to be able to send requests. 

In this first version of the application `fmtok8s-app` all the interactions between the services happen via REST calls. This push the caller to know the other services names. 

You can open different tabs in Cloud Shell to inspect the logs of each service when you are using the application (submitting and approving/rejecting proposals). 


## Challenges
This section covers some of the challenges that you might face when working with these kind of applications inside Kubernetes. This section is not needed to continue with the workshop, but it highlight the need for some other tools to be used in conjuction with the application. 
 
<details>
  <summary>To see more details about the challenges Click to Expand</summary>

Among some of the challenges that you might face are the following big topics:
- **Flow buried in code and visbility for non-technical users**: for this scenario the `C4P` service is hosting the core business logic on how to handle new proposals. If you need to explain to non-technical people how the flow goes, you will need to dig in the code to be 100% sure about what the application is doing. Non-technical users will never sure about how their applications are working, as they can only have limited visibility of what is going on. How would you solve the visibility challenge? 
- **Edge Cases and Errors**: This simple example, shows what is usually called the **happy path**, where all things goes as expected. In real life implementations, the amount of steps that happens to deal with complex scenarios grows. If you want to cover and have visibility on all possible edge cases and how the organization deals with logical errors that might happen in real life situations, it is key to document and expose this information not only to technical users. How would you document and keep track of every possible edge case and errors that can happen in your distributed applications? How would you do that for a monolith application?
- **Dealing with changes**: for an organization, being able to understand how their applications are working today, compared on how they were working yesterday is vital for communication, in some cases for compliance and of course to make sure that different deparments are in sync. The faster that you want to go with microservices, the more you need to look outside the development departments to apply changes into the production environments. You will need tools to make sure that everyone is in the same page when changes are introduced. How do you deal with changes today in your applications? How do you expose to non-technical users the differences between the current version of your application that is running in production compared with the new one that you want to promote?
- **Implementing Time-Based Actions/Notifications**: I dare to say that about 99% of applications require some kind of notification mechanism that needs to deal with scheduled actions at some point in the future or require to be triggered every X amount of time. When working with distributed system, this is painful, as you will need a distributed scheduler to guarantee that things scheduled to trigger are actually triggered when the time is right. How would you implement these time based behaviours? If you are thinking about **Kubernetes Cron Jobs** that is definitely a wrong answer. 
- **Reporting and Analytics**: if you look at how the services of applications are storing data, you will find out that the structures used are not optimized for reporting or doing analytics. A common approach, is to push the infomration that is relevant to create reports or to do analytics to [ElasticSearch](https://www.elastic.co/elastic-stack) where data can be indexed and structured for efficient querying. Are you using ElasticSearch or a similar solution for building reports or running analytics? 

In version 2 of the application you will be working to make your application's internals more visibile to non-technical users. 

</details> 

You will now undeploy version 1 of the application to deploy version 2. You only need to undeploy version 1 to save resources.
In order to undeploy version 1 of the application you can run:
``` bash
h delete fmtok8s --no-hooks
```


## Version 2: Knative, Cloud Events and Camunda Cloud

Version 2 of the application is configured to emit [Cloud Events](http://cloudevents.io), whenever something relevant happens in any of the services. For this example, you are interested in the following events: 
- `Proposal Received`
- `Proposal Decision Made`
- `Email Sent`
- In the case of the proposal being approved `Agenda Item Created` 

Version 2 of the application still uses the same version of the services found in Version 1, but these services are configured to emit events to a **Knative Broker** that was created when you installed Knative. This Knative Broker, receive events and routed them to whoever is interested in them. In order to register interest in certain events, Knative allows you to create **Triggers** (which are like subscriptions with filters) for this events and specify where these events should be sent. 

For Version 2, you will use the **Zeebe Workflow Engine** provisioned in your **Camunda Cloud** account to capture and visualize these meaninful events.
In order to route these **Cloud Events** from the Knative Broker to **Camunda Cloud** a new component is introduced along your Application services. This new component is called **Zeebe Cloud Events Router** and serves as the bridge between Knative and Camunda Cloud, using Cloud Events as the standardize communication protocol. 

As you can imagine, in order for the **Zeebe Cloud Events Router** to connect with your **Camunda Cloud Zeebe Cluster** you need to create a new **Client**, a set of credentials which allows these components to connect and communicate. 

Go to the **Camunda Cloud** console, click on your cluster to see your cluster details:

<img src="workshop-imgs/19-cluster-details.png" alt="Cluster Details" width="500px">

Go to the **Clients** tab and then **Create a New Client**:

<img src="workshop-imgs/21-create-cluster-client.png" alt="Cluster Details" width="500px">

Call it `my-client` and click **Add**:

<img src="workshop-imgs/22-call-it-my-client.png" alt="Cluster Details" width="500px">

The new client called `my-client` will be created: 

<img src="workshop-imgs/23-client-created.png" alt="Cluster Details" width="500px">

Now you can access the **Connection Information**:

<img src="workshop-imgs/24-client-information-kube-secret.png" alt="Cluster Details" width="500px">

By clicking the button **Copy Kubernetes Secret** the command will be copied into your clipboard and you can paste it inside **Cloud Shell** inside Google Cloud.

``` bash
k create secret generic camunda-cloud-secret --from-literal=ZEEBE_ADDRESS=...
```

By running the previous command, you have created a new `Kubernetes Secret` that host the credentials for our applications to talk to Camunda Cloud. 
Now you are ready to install version 2 of the application by running (again ignore the warnings): 

``` bash 
h install fmtok8s-v2 workshop/fmtok8s-app-v2
```

<img src="workshop-imgs/29-helm-install-v2.png" alt="Cluster Details" width="700px">

You can check that all the services are up and running with the same two commands as before:

``` bash
k get pods
```

and

``` bash
k get ksvc
```

You should see something like this:

<img src="workshop-imgs/30-k-get-pod-and-ksvc.png" alt="Cluster Details" width="700px">

Notice that now the **Zeebe Cloud Events Router** is running along side the application services, and it is configured to use the Kubernetes Secret that was previously created to connect to **Camunda Cloud**.

But here is still one missing piece to route the **Cloud Events** generated by your application services to the **Zeebe Cloud Events Router** and those are the Knative Triggers (Subscriptions to route the events from the broker to wherever you want). 

These Knative Triggers are defined in YAML and can be packaged inside the Application V2 Helm Chart, which means that they are installed as part of the application. 

You can list these Knative Triggers by running the following command:
``` bash
k get triggers
```

You should see an output like this: 

<img src="workshop-imgs/31-k-get-triggers.png" alt="Cluster Details" width="700px">

Finally, even when Cloud Events are being routed to Camunda Cloud, you need to create a model that will consume the events that are coming from the application, so they can be correlated and visualized. 

You can download the models that you will be using in [the next steps from here](https://github.com/salaboy/from-monolith-to-k8s-assets/archive/1.0.0.zip).

Once you downloaded the models, extract the ZIP file a place that you can quickly locate to upload these files in the next steps. 

Now, go back to your **Camunda Cloud Zeebe Cluster** list (you can do this by clicking in the top breadcrum with the name of your Organization):

<img src="workshop-imgs/32-camunda-cloud-cluster-list.png" alt="Cluster Details" width="700px">

Next, click on the **BPMN Diagrams(beta)** Tab, then click **Create New Diagram**:

<img src="workshop-imgs/33-bpmn-diagrams-list.png" alt="Cluster Details" width="700px">

With the Diagram editor opened, first enter the name **visualize** into the diagram name box and then click the **Import Diagram** button:

<img src="workshop-imgs/34-name-and-import-diagram.png" alt="Cluster Details" width="700px">

Now choose **c4p-visualize.bpmn** from your filesystem: 

<img src="workshop-imgs/35-choose-c4p-visualize-from-filesystem.png" alt="Cluster Details" width="700px">

The diagram shoud look like:

<img src="workshop-imgs/36-diagram-1-should-look-like.png" alt="Cluster Details" width="700px">

With the Diagram ready, you can now hit **Save and Deploy**:

<img src="workshop-imgs/37-save-and-deploy.png" alt="Cluster Details" width="700px">

Next, **close/disregard** the popup suggesting to start a new instance:

<img src="workshop-imgs/38-close-popup.png" alt="Cluster Details" width="700px">

Well Done! you made it, now everything is setup for routing and fowarding events from our application, to Knative Eventing, to the Zeebe Cloud Events Router to Camunda Cloud. 


In order to see how this is actually working you can use **Camunda Operate**, a dashboard included inside **Camunda Cloud** which allows you to understand how these models are being executed, where things are at a giving time and to troubleshoot errors that might arise from your applications daily operations.

You can access **Camunda Operate** from your cluster details, inside the **Overview Tab**, at the bottom, clicking in the **View in Operate** link:

<img src="workshop-imgs/39-cluster-details-operate-link.png" alt="Cluster Details" width="700px">

You should see the **Camunda Operate** main screen, where you can click in the **C4P Visualize** section highlighted in the screenshot below:

<img src="workshop-imgs/40-operate-main-screen.png" alt="Cluster Details" width="700px">

This opens the runtime data associated with our workflow models, now you should see this:

<img src="workshop-imgs/41-visualize-diagram-in-operate.png" alt="Cluster Details" width="700px">

Now go back to the Conference Application, remember, listing all the Knative Services will show the URL for the API Gateway Service that is hosting the User Interface, when you are in the application, submit a new proposal and then refresh **Camunda Operate**:

<img src="workshop-imgs/42-runtime-data-in-operate.png" alt="Cluster Details" width="700px">

If you click into the Instance ID link, highligted above, you can see the details of that specific instance:

<img src="workshop-imgs/43-instance-data.png" alt="Cluster Details" width="700px">

If you go ahead to the **Back Office** of the application and **approve** the proposal that you just submitted, you should see in **Camunda Operate** that the instance is completed:

<img src="workshop-imgs/44-approved-instance.png" alt="Cluster Details" width="700px">

<details>
  <summary>+ Extras: Understanding different branches (Click to Expand)</summary>
As you might notice, the previous workflow model will only work if you approve proposals, as the `Agenda Item Created` event is only emitted if the proposal is accepted. In order to cover also the case when you reject a proposal you can deploy version 2 of the workflow model, that describes these two branches for approving and rejecting proposals.
  
In order to deploy the second version of the workflow model you follow the same steps as before, you   
  
  <img src="workshop-imgs/45-bpmn-diagrams-list-with-v1.png" alt="Cluster Details" width="700px">

Click into your previously saved diagram called **visualize** and then **Import Diagram** and then select **cp4-visualize-with-branches.bpmn**:

<img src="workshop-imgs/46-c4p-visualize-v2.png" alt="Cluster Details" width="700px">

The new diagram should like this:

<img src="workshop-imgs/47-c4p-vizualize-diagram.png" alt="Cluster Details" width="700px">

Now you are ready to **Save and Deploy** the new version:

<img src="workshop-imgs/48-save-and-deploy-v2.png" alt="Cluster Details" width="700px">

If you switch back to **Camunda Operate** you will now see two versions of the **C4P Visualize** workflow:

<img src="workshop-imgs/49-camunda-operate-two-versions.png" alt="Cluster Details" width="700px">

Click to drill down into the runtime data for the new version of the workflow:

<img src="workshop-imgs/50-new-version-deployed.png" alt="Cluster Details" width="700px">

If you now go back to the application and submit two proposals you can reject one and approve one, you should now see both instances completed:

<img src="workshop-imgs/51-version-2-completed.png" alt="Cluster Details" width="700px">

Remember that you can click in any instance to find more details about the execution, such as the audit logs to understand exactly when things happened:

<img src="workshop-imgs/52-audit-log.png" alt="Cluster Details" width="700px">

</details>

Let's undeploy version 2 to make some space for version 3. 

```
h delete fmtok8s-v2 --no-hooks
```

# Version 3; Workflow Orchestration with Camunda Cloud

In Version 3, you will orchestrate the interactions using the workflow engine. You can now install version 3 running:

``` bash
h install fmtok8s-v3 workshop/fmtok8s-app-v3
```
The ouput for this command should look familiar at this stage:

<img src="workshop-imgs/53-installing-v3.png" alt="Cluster Details" width="700px">

Check that the Kubernetes Pods and the Knative Services are ready:

<img src="workshop-imgs/54-checking-pods-ksvc-v3.png" alt="Cluster Details" width="700px">

When all the pods are ready (2/2) you can now access to the application. 

As you might have noticed, there is a new Knative Service and pod called **fmtok8s-speakers**, you will use that service later on.  

An important change in version 3 is that it doesn't use a REST based communication between services, this version let **[Zeebe](http://zeebe.io)**, the workflow engine inside **Camunda Cloud**, to define the sequence and orchestrate the services interactions. **[Zeebe](http://zeebe.io)** uses a Pub/Sub mechanism to communicate with each service, which introduces automatic retries in case of failure and reporting incidents when there are service failures. 

<details>
  <summary>+ Extras: Changes required to let Zeebe communicate with our existing services (Click to Expand)</summary>
Links to Workers, and dependencies in projects, plus explain how the workers code is reusing the same code as rest endpoints internally. 
</details>

Another important change, is that the **C4P Service** now deploys automatically the workflow model used for the orchestration to happen. 
This means that when the **fmtok8s-c4p** Knative Serivce is up and ready, you should have a new workflow model already deployed in **Camunda Cloud**:

<img src="workshop-imgs/55-new-workflow-model-in-v3.png" alt="Cluster Details" width="700px">

If you now click into the new workflow model you can see how the new model looks like: 

<img src="workshop-imgs/56-v3-orchestration-workflow-model.png" alt="Cluster Details" width="700px">

If you submit a **new proposal** from the application user interface, this new workflow model is in charge of defining the order in which services are invoked. 
From the end user point of view, nothing has changed, besides the fact that they can now use **Camunda Operate** to understand in which step each proposal is at a given time. From the code perspective, the business logic required to define the steps is now delegated to the workflow engine, which enables non-technical people to gather valuable data about how the organization is working, where the bottlenecks are and how are your Cloud-Native applications working. 

Having in a single place the state of all proposals can help organizers to prioritize other work or just make decisions to move things forward:

<img src="workshop-imgs/57-quick-overview-of-state.png" alt="Cluster Details" width="700px">

In the screenshot above, it is clear that 2 proposals are still waiting for a decision, 2 proposals were approved and 1 rejected. 
Remember that you can drill down to each individual workflow instance for more details, for example, how much time a proposal has been waiting for a decision:

<img src="workshop-imgs/58-waiting-for-decision.png" alt="Cluster Details" width="700px">

Based on the data that the workflow engine is collecting from the workflow's executions, you can understand better where the bottlenecks are or if there are simple things that can be done to improve how the organization is dealing with proposals. For this example, you can say that this workflow model represent 100% the steps required to accept or reject a proposal, in some way this explains to non-technical people the steps that the application is executing under the covers. 

Becuase the workflow model is now in charge of the sequence of interactions, you are free to change and adapt the workflow model to better suit your organization needs. 

If you made it this far, Well Done!!!

Here are some extras that you might be interested in to make more sense of what you have done so far:
- Update the workflow model to use the newly introduced Speakers Service
- Update the workflow model to send notifications if a proposal is waiting for a decision for too long
- Make the application fail to see how incidents are reported into Camunda Operate

# Next Steps



# Thanks to all contributors
- [MCruzDev1]()
- [Ray Tsang]()
- 

