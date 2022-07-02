# Identity Management and Single Sign On

This short document explains the changes to configure the Conference Platform to support Single Sign On and Role Based Access Control to the Backoffice, 
where only Users with the Role Organizers should be able to access. 

There are two very complete solutions that you should look at in this space, both Open Source, but Zitadel also provides a SaaS approach: 

- Zitadel
- Keycloak 

## Updating the conference application to use Zitadel for SSO and Identity Management

If we are looking at Zitadel, the first thing that we need to do is to decide if we want to use the Hosted/Managed version or if we want to install Zitadel in our Kubernetes Cluster. 
This tutorial will go over installing Zitadel into the Kubernetes cluster where we have our application running. 

To install Zitadel you can use a Helm Chart or you can install Zitadel as a Knative Service. 

Links: 
- React: https://docs.zitadel.com/docs/quickstarts/login/react 
- Spring Boot Example: https://github.com/zitadel/zitadel-examples/blob/main/java/spring-boot/api/src/main/resources/application.yml
- Helm chart installation and Knative: 
