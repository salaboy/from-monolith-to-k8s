# JBCNConf 2022 Script (Spanish)

Agenda
- [Intro / Background](#intro--background)
- [IDEs, Lenguajes y Frameworks]()
- [Hablemos de Containers y Kubernetes]()
- [Extendiendo Kubernetes]()
- [Alternativas]()

## Intro / Background

[@Salaboy](https://twitter.com/salaboy) / [https://salaboy.com](https://salaboy.com)
![avatar.png](avatar.png)
- Java 
- J2EE
- JBoss 
  -  Java EE
  -  Wildfly
  -  Kubernetes
- Mi primer Kubernetes Controller con Fabric8.io
- Jenkins X
- Spring Boot y Spring Cloud 
  - Spring Cloud Kubernetes
  - Kubernetes Controllers con Spring Cloud Kubernetes
- Go
  - Kubernetes Controllers con KubeBuilder
  - Knative
   - Knative Eventing 
   - Knative Functions Working Group co-lead
- [Continuous Delivery for Kubernetes]()
![book](book.png)

## IDEs, Lenguajes y Frameworks

En el contexto de crear un servicio que expone un endpoint REST.

Goland and Intellij Idea nos hacen la vida muy facil ya que la experiencia es la misma.

Vamos a ver codigo. Empezamos por Spring Boot y Quarkus, pero ya que es una conferencia de Java no vamos a entrar en detalles. 

Las aplicaciones que voy a mostrar esta en este repository: 
- [Spring Boot](spring-boot/conference-service/)
- [Quarkus](quarkus/conference-service/)
- [Go](go/conference-service/)

En resumen: 
- Go soluciona depedency management (Go Modules) sin tener que usar una herramienta externa
- Go require mas conocimiento del ecosistema a la hora de escoger librarias (usamos Gorilla Mux, pero hay muchas mas)
- Go tiene incluido marshallers para JSON y YAML
- Go crea binarios que dependenden de la plataform donde hagamos el build

## Hablemos de Containers y Kubernetes

Como vamos de los proyectos que vimos antes a tener containers creados 


