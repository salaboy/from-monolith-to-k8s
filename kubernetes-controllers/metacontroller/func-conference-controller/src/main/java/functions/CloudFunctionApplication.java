package functions;

import java.util.*;
import java.util.function.Function;

import com.fasterxml.jackson.annotation.JsonProperty;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.annotation.Bean;
import org.springframework.http.MediaType;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.reactive.function.client.WebClient;
import org.yaml.snakeyaml.Yaml;
import reactor.core.publisher.Mono;

@SpringBootApplication
public class CloudFunctionApplication {

  private static final Logger log = LoggerFactory.getLogger(CloudFunctionApplication.class);

  public static void main(String[] args) {
    SpringApplication.run(CloudFunctionApplication.class, args);
  }

//  @Autowired
//  private WebClient.Builder webClient;

  @Bean
  WebClient webClient(WebClient.Builder b) {
    return b.build();
  }

  @Bean
  public Function<ControllerInput, Mono<DesiredState>> reconcile(WebClient http) {
    // Input: A controller JSON object that contains a CompositeController, 
    // a Parent and Child, where the Parent describes the user input that triggered the function call
    // - The CompositeController is controller responsible for the webhook
    // - The parent is the triggering input - the resource that needs to be reconciled
    // - The children are the resources that can be modified as part of the reconciliation loop
    // See CompositeController.Spec for a summary of the triggering resource, the hook, and the children (the resources that can be changed)
    return (resource) -> {
      Parent parent = resource.parent();
      ParentMetadata parentMetadata = parent.metadata();
      ParentSpec parentSpec = parent.spec();
      
      log.info("Reconciling Resource: " + parent.apiVersion() + "/" + parent.kind() + " > " + parentMetadata.name());

      boolean productionTestEnabled = parentSpec.productionTestEnabled();
      String conferenceNamespace = parentSpec.namespace();
      
      List<Child> children = new ArrayList<Child>();
      
      return Mono.zip(
          getServiceInfo(http, "fmtok8s-frontend", conferenceNamespace),
          getServiceInfo(http, "fmtok8s-email", conferenceNamespace),
          getServiceInfo(http, "fmtok8s-agenda", conferenceNamespace),
          getServiceInfo(http, "fmtok8s-c4p", conferenceNamespace)
      ).map(serviceInfos -> {
          log.info("Service Infos: " + serviceInfos);
          boolean frontendReady = false;
          boolean agendaServiceReady = false;
          boolean emailServiceReady = false;
          boolean c4pServiceReady = false;
          if (!serviceInfos.getT1().contains("N/A") && !serviceInfos.getT1().isEmpty()) {
            frontendReady = true;
          }
          if (!serviceInfos.getT2().contains("N/A") && !serviceInfos.getT2().isEmpty()) {
            emailServiceReady = true;
          }
          if (!serviceInfos.getT3().contains("N/A") && !serviceInfos.getT3().isEmpty()) {
            agendaServiceReady = true;
          }
          if (!serviceInfos.getT4().contains("N/A") && !serviceInfos.getT4().isEmpty()) {
            c4pServiceReady = true;
          }
          
          boolean conferenceReady = false;
          if (frontendReady && emailServiceReady && agendaServiceReady && c4pServiceReady) {
            conferenceReady = true;
            if (productionTestEnabled) {
              Map<String, Object> deployment = createProductionTestDeployment();
              children.add(new Child(deployment));
            }
          }
          
          String url = "Impossible to know without access to the K8s API";
          Status status = new Status(frontendReady, emailServiceReady, agendaServiceReady, c4pServiceReady, productionTestEnabled, conferenceReady, url);

          DesiredState desiredState = new DesiredState(children, status);

          log.info("> Desired State: " + desiredState);
          return desiredState;
        });
    };
  }
  
//  public Mono<String> getServiceInfo(String url) {
//    return webClient.build()
//      .get()
//      .uri(url)
//      .accept(MediaType.APPLICATION_JSON)
//      .retrieve()
//      .bodyToMono(String.class)
//      .onErrorResume(err -> Mono.just("N/A"));
//
//  }

  public Mono<String> getServiceInfo(WebClient http, String name, String namespace) {
    return getServiceInfo(http, "http://" + name + "." + namespace + ".svc.cluster.local/info");
  }
  
  public Mono<String> getServiceInfo(WebClient http, String url) {
    return http
      .get()
      .uri(url)
      .accept(MediaType.APPLICATION_JSON)
      .retrieve()
      .bodyToMono(String.class)
      .onErrorResume(err -> Mono.just("N/A"));

  }

  public Map<String, Object> createProductionTestDeployment() {
    Yaml yaml = new Yaml();
    String deploymentYaml = "apiVersion: apps/v1\n" +
      "kind: Deployment\n" +
      "metadata:\n" +
      "  name: metacontroller-production-tests\n" +
      "spec:\n" +
      "  replicas: 1\n" +
      "  selector:\n" +
      "    matchLabels:\n" +
      "      app: production-tests\n" +
      "  template:\n" +
      "    metadata:\n" +
      "      labels:\n" +
      "        app: production-tests\n" +
      "    spec:\n" +
      "      containers:\n" +
      "        - name: production-tests\n" +
      "          image: salaboy/metacontroller-production-tests:metacontroller\n" +
      "          imagePullPolicy: Always\n";
    return yaml.load(deploymentYaml);
  }

  // ControllerInput: A controller JSON object that contains a CompositeController, 
  // a Parent and Child, where the Parent describes the user input that triggered the function call
  // 1. The CompositeController is controller responsible for the webhook
  // The parent is the triggering input - the resource that needs to be reconciled
  // The children are the resources that can be modified as part of the reconciliation loop
  // See CompositeController.Spec for a summary of the triggering resource, the hook, and the children (the resources that can be changed)
  
  record ParentMetadata(String name){}

  record ParentSpec(String namespace, @JsonProperty("production-test-enabled") boolean productionTestEnabled) {
  }

  record Parent(String apiVersion, String kind, @JsonProperty("metadata") ParentMetadata metadata,
                @JsonProperty("spec") ParentSpec spec) {
  }

  // Input object. Only Parent (and Children?) are needed for the purposes of the demo
  record ControllerInput(Parent parent) {
  }

  record Child(@JsonProperty("Deployment.apps/v1") Map<String, Object> productionTestDeployment){}

  record Status(
    @JsonProperty("frontend-ready") boolean frontendReady, 
    @JsonProperty("email-service-ready") boolean emailServiceReady, 
    @JsonProperty("agenda-service-ready") boolean agendaServiceReady, 
    @JsonProperty("c4p-service-ready") boolean c4pServiceReady, 
    @JsonProperty("prod-tests") boolean productionTestEnabled,
    @JsonProperty("ready") boolean conferenceReady,
    String url){
  }

  // Output object. Should contain desired state with children and status
  record DesiredState(List<Child> children, Status status) {
  }
  
}
