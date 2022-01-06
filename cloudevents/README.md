# CloudEvents Tutorial

In this short tutorial we are going to create two applications that produce and consume CloudEvents. These applications uses different technology stacks: 
- Application A (`fmtok8s-java-cloudevents`) uses Java and Spring Boot and it adds the CloudEvents Java SDK to write and read CloudEvents. 
- Application B (`fmtok8s-go-cloudevents`) uses Go and adds the CloudEvents Go SDK to read and write CloudEvents. 

![CloudEvents Examples](cloudevents-fmtok8s.png)

If you want to build and run the applications you will need to hava Java, Maven and Go installed. Alternatively you can run the available Docker containers, in which case you only need Docker to run this examples. 


To consume CloudEvents via HTTP both applications expose a REST endpoint where the events can be received. 

