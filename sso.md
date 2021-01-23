# From Monolith to K8s - Workshop 


## Installing and Configuring Keycloak
During this step-by-step you will be using **Kubernetes Cluster** and a Keycloak as SSO to secure our API Gateway and Microservices. 

### Creating a Kubernetes Cluster with KIND

```
$ kind create cluster --name keycloak
```

Don't forget to set current cluster/context

```
$ kubectl cluster-info --context kind-keycloak
```

In this example we'll not create a namespace (It is not a best practice )

### Adding Keycloak on Cluster Kubernetes

```
kubectl cretate -f https://raw.githubusercontent.com/keycloak/keycloak-quickstarts/latest/kubernetes-examples/keycloak.yaml
```

Let's see the keycloak pod

```
$ kubectl get pods
```

## Configuring Keycloak
### 1 - Let's access Administration Console:

<img src="sso-imgs/sso-1.png" alt="Go to Administration Console" width="700px">

### 2 - We'll access using our credentials passed through configurations

<img src="sso-imgs/sso-2.png" alt="Login" width="700px">

### 3 - Let's create our realm (fmtok8s)

<img src="sso-imgs/sso-3.png" alt="Creating Keycloak Realm" width="700px">

### 4 - Creating a Realm's Client

<img style="margin-bottom: 10px;" src="sso-imgs/sso-4.png" alt="Creating Realm's Client part 1" width="700px">

<img src="sso-imgs/sso-5.png" alt="Creating Realm's Client part 2" width="700px">

### 5 - Configuring a Client

The client configuration's page is very large, then I will divide it in two parts:

Parte 1:
<img style="margin-bottom: 10px;" src="sso-imgs/sso-6.png" alt="Configuring The Client Gateway part 1" width="700px">

Parte 2:
<img src="sso-imgs/sso-7.png" alt="Configuring The Client Gateway part 2" width="700px">

## 6 - Creating an Realm's User

<img src="sso-imgs/sso-8.png" alt="Creating an user to Realm" width="700px">

<img src="sso-imgs/sso-9.png" alt="Adding an user to Realm" width="700px">

After, you should set the user's password

<img src="sso-imgs/sso-10.png" alt="Setting user's password" width="700px">

### 7 - Creating a Realm's Role

<img src="sso-imgs/sso-11.png" alt="Creating a Realm's Role" width="700px">

<img src="sso-imgs/sso-12.png" alt="Adding Realm's Role" width="700px">

### 8 - Adding a role to user

<img src="sso-imgs/sso-13.png" style="margin-bottom: 20px;" alt="Adding Realm's Role to User" width="700px">


## Changing API Gateway to secure our hidden microservices

[API Gateway](https://github.com/mcruzdev/fmtok8s-api-gateway) was created with Spring Cloud Gateway. The Spring Cloud Gateway uses Spring Webflux working with reactive stack.

There is a great lib called
`org.keycloak:keycloak-spring-boot-starter` that help us to configure our application using keycloak and it runs better with Servlet applications. [See](https://keycloak.discourse.group/t/webflux-support-for-spring-boot-and-spring-security-adapters/2936)

In this workshop, you will use Spring Security OAuth2. Let's go to use it.

### Adding Spring OAuth2 dependecies in API Gateway

```
<dependency>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-starter-oauth2-client</artifactId>
</dependency>
<dependency>
    <groupId>org.springframework.cloud</groupId>
    <artifactId>spring-cloud-starter-security</artifactId>
</dependency>
```

### Configuring 

We should change the configuration of **API Gateway** on application.yml

```
spring:
  security:
    oauth2:
      client:
        provider:
          oidc:
            issuer-uri: http://localhost:8081/auth/realms/fmtok8s
        registration:
          oidc:
            client-name: keycloak
            provider: oidc
            client-id: gateway
            client-secret: 7208a758-e57c-4045-8c4a-9831bb2b8144
```

The property client-secret should be filled with client-secret from keycloak, let's get our client-secret from gateway keycloak client:

**Copy gateway's secret:**

<img src="sso-imgs/sso-14.png" alt="Adding Realm's Role" width="700px">


And add a filter a new filter on **API Gateway** to relay de Token to hiden services

```
  cloud:
    gateway:
      default-filters:
        - TokenRelay=
        - RemoveRequestHeader=Cookie
      httpclient:      
``` 

### Creating class SecurityConfig to our Gateway

```
package com.salaboy.conferences.site.security;

import org.springframework.context.annotation.Bean;
import org.springframework.security.config.annotation.web.reactive.EnableWebFluxSecurity;
import org.springframework.security.config.web.server.ServerHttpSecurity;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.authority.SimpleGrantedAuthority;
import org.springframework.security.oauth2.client.oidc.userinfo.OidcReactiveOAuth2UserService;
import org.springframework.security.oauth2.client.oidc.userinfo.OidcUserRequest;
import org.springframework.security.oauth2.client.userinfo.ReactiveOAuth2UserService;
import org.springframework.security.oauth2.core.oidc.user.DefaultOidcUser;
import org.springframework.security.oauth2.core.oidc.user.OidcUser;
import org.springframework.security.oauth2.core.oidc.user.OidcUserAuthority;
import org.springframework.security.web.server.SecurityWebFilterChain;

import java.util.*;
import java.util.stream.Collectors;

@EnableWebFluxSecurity
public class SecurityConfig {

    @Bean
    public SecurityWebFilterChain springSecurityFilterChain(ServerHttpSecurity http) {
        return http.csrf().disable()
                .authorizeExchange()
                .pathMatchers("/backoffice/**").hasRole("approver")
                .anyExchange().permitAll()
                .and()
                .oauth2Login()
                .and()
                .oauth2Client()
                .and()
                .build();
    }

    @Bean
    public ReactiveOAuth2UserService<OidcUserRequest, OidcUser> oidcUserService() {
        final OidcReactiveOAuth2UserService delegate = new OidcReactiveOAuth2UserService();

        return (userRequest) -> {
            // Delegate to the default implementation for loading a user
            return delegate.loadUser(userRequest).map(user -> {
                Set<GrantedAuthority> mappedAuthorities = new HashSet<>();

                user.getAuthorities().forEach(authority -> {
                    if (authority instanceof OidcUserAuthority) {
                        OidcUserAuthority oidcUserAuthority = (OidcUserAuthority) authority;
                        mappedAuthorities.addAll(extractAuthorityFromClaims(oidcUserAuthority.getUserInfo().getClaims()));
                    }
                });

                return new DefaultOidcUser(mappedAuthorities, user.getIdToken(), user.getUserInfo());
            });
        };
    }

    public static List<GrantedAuthority> extractAuthorityFromClaims(Map<String, Object> claims) {
        return mapRolesToGrantedAuthorities(getRolesFromClaims(claims));
    }

    @SuppressWarnings("unchecked")
    private static Collection<String> getRolesFromClaims(Map<String, Object> claims) {
        return (Collection<String>) claims.getOrDefault("groups",
                claims.getOrDefault("roles", new ArrayList<>()));
    }

    private static List<GrantedAuthority> mapRolesToGrantedAuthorities(Collection<String> roles) {
        return roles.stream()
                .map("ROLE_"::concat)
                .map(SimpleGrantedAuthority::new)
                .collect(Collectors.toList());
    }
}
```


## Securing our microservices (C4P)

### We need to add some dependecies on pom.xml

```
        <dependency>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-starter-security</artifactId>
        </dependency>

        <dependency>
            <groupId>org.springframework.security</groupId>
            <artifactId>spring-security-oauth2-resource-server</artifactId>
        </dependency>

        <dependency>
            <groupId>org.springframework.security</groupId>
            <artifactId>spring-security-oauth2-jose</artifactId>
        </dependency>
```

### Configuring Spring OAuth2 Resource Server through application.properties:

```
spring.security.oauth2.resourceserver.jwt.issuer-uri=http://localhost:8080/auth/realms/fmtok8s
```

### Creat an class to configure CORS 

```
package com.salaboy.conferences.c4p.rest.configuration;

import org.springframework.context.annotation.Configuration;
import org.springframework.web.reactive.config.CorsRegistry;
import org.springframework.web.reactive.config.WebFluxConfigurer;

@Configuration
public class CORSConfig implements WebFluxConfigurer {

    @Override
    public void addCorsMappings(CorsRegistry registry) {
        registry.addMapping("/**").allowCredentials(true).allowedMethods("*");
    }
}
```

### Creating our SecurityConfig

```
package com.salaboy.conferences.c4p.rest.configuration;

import org.springframework.context.annotation.Bean;
import org.springframework.core.convert.converter.Converter;
import org.springframework.http.HttpMethod;
import org.springframework.security.authentication.AbstractAuthenticationToken;
import org.springframework.security.config.annotation.web.reactive.EnableWebFluxSecurity;
import org.springframework.security.config.web.server.ServerHttpSecurity;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.authority.SimpleGrantedAuthority;
import org.springframework.security.oauth2.jwt.Jwt;
import org.springframework.security.oauth2.server.resource.authentication.JwtAuthenticationConverter;
import org.springframework.security.oauth2.server.resource.authentication.ReactiveJwtAuthenticationConverterAdapter;
import org.springframework.security.web.server.SecurityWebFilterChain;
import reactor.core.publisher.Mono;

import java.util.Collection;
import java.util.Collections;
import java.util.List;
import java.util.stream.Collectors;

@EnableWebFluxSecurity
public class SecurityConfig {

    @Bean
    public SecurityWebFilterChain securityWebFilterChain(ServerHttpSecurity http) {

        http
        .csrf().disable()
        .authorizeExchange(exchanges ->
                exchanges
                        .pathMatchers(HttpMethod.POST, "/**").hasAnyAuthority("approver")
                        .pathMatchers(HttpMethod.DELETE, "/**").hasAnyAuthority("approver")
                        .anyExchange().permitAll())
        .oauth2ResourceServer(oauth2 ->
                oauth2.jwt(jwt -> jwt.jwtAuthenticationConverter(grantedAuthoritiesExtractor())));

        return http.build();
    }

    Converter<Jwt, Mono<AbstractAuthenticationToken>> grantedAuthoritiesExtractor() {
        JwtAuthenticationConverter jwtAuthenticationConverter =
                new JwtAuthenticationConverter();

        jwtAuthenticationConverter.setJwtGrantedAuthoritiesConverter(new GrantedAuthoritiesExtractor());

        return new ReactiveJwtAuthenticationConverterAdapter(jwtAuthenticationConverter);
    }

    static class GrantedAuthoritiesExtractor implements Converter<Jwt, Collection<GrantedAuthority>> {

        @Override
        public Collection<GrantedAuthority> convert(Jwt jwt) {

            @SuppressWarnings("unchecked")
            var roles = (List<String>) jwt.getClaims().getOrDefault("groups", Collections.emptyList());

            return roles.stream()
                    .map(SimpleGrantedAuthority::new)
                    .collect(Collectors.toList());
        }
    }
}

```