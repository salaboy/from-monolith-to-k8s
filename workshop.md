# Workshop 

During this workshop you will deploy a Cloud Native application, inspect it, change its configuration to use different services and 
play around with it to get familiar with Kubernetes and Cloud Native tools that can help you to be more efficient. 

During this workshop you will be using GKE (Managed Kubernetes Engine inside Google Cloud) to deploy a complex application composed by multiple services. In order to do this 
you will be using `kubectl` and `helm` to deploy the application. Because you will be using the Google Cloud Console, you can save some time by creating some aliases for these two commands

```
> alias k=kubectl
> alias h=helm
```
Now instead of typing `kubectl` or `helm` you can just type `k` and `h` respectivily. 

Once you have these alias set up you can proceed to add a new Helm Repository where the Helm packages for the application are stored. 
You can do this by runnig the following command

```
> h repo add workshop http://chartmuseum-jx.35.222.17.41.nip.io
> h repo update
```
