# JBCNConf 2022 Script (Spanish)

Agenda
- [Intro / Background](#intro--background)
- [IDEs, Lenguajes y Frameworks](#ides-lenguajes-y-frameworks)
- [Hablemos de Containers y Kubernetes](#hablemos-de-containers-y-kubernetes)
- [Extendiendo Kubernetes](#extendiendo-kubernetes)
- [Alternativas más saludables](#alternativas-mas-saludables)

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

Como hacemos para que los proyectos que creamos corran dentro de Kubernetes? 
Vamos de vuelta a los proyectos: 

- [Spring Boot](spring-boot/conference-service/)
- [Quarkus](quarkus/conference-service/)
- [Go](go/conference-service/)

En resumen: 
- Spring Boot y Quarkus proveen integraciones con Jib y Buildpacks para contruir contenedores sin tener que definir Dockerfiles
  - Ambas integraciones usan la version definida en Maven para taggear el container
- En Go podemos usar Ko para construir y publicar estos containers a nuestro registry preferido
- Spring Boot out of the box no provee generacion de YAMLs, se puede usar JKube para esto
- Quarkus provee generacion de YAMLs con una extension
- Tanto en Spring Boot y en Quarkus nosotros tenemos que publicar nuestros containers
- `ko` publica los containers y usa SHA en vez de una version fija. Esto nos permite construir y correr containers con los ultimos cambios
  - Con `ko` tenemos que crear nuestros propios YAMLs pero `ko resolve` reemplaza las referencias a los builds de los containers


## Extendiendo Kubernetes



## Alternativas más saludables

