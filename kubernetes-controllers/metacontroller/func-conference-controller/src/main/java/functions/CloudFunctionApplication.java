package functions;

import domain.controllerinput.ControllerInput;
import domain.controllerinput.Parent;
import domain.controllerinput.ParentMetadata;
import domain.controllerinput.ParentSpec;
import domain.controlleroutput.ControllerOutput;
import domain.controlleroutput.Status;
import io.kubernetes.client.common.KubernetesObject;
import io.kubernetes.client.openapi.models.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.annotation.Bean;
import org.springframework.http.MediaType;
import org.springframework.web.reactive.function.client.WebClient;
import reactor.core.publisher.Mono;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.function.Function;

@SpringBootApplication
public class CloudFunctionApplication {

  private static final Logger log = LoggerFactory.getLogger(CloudFunctionApplication.class);

  public static void main(String[] args) {
    SpringApplication.run(CloudFunctionApplication.class, args);
  }

  @Bean
  WebClient webClient(WebClient.Builder b) {
    return b.build();
  }

  @Bean
  public Function<ControllerInput, Mono<ControllerOutput>> reconcile(WebClient http) {
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
      
      List<KubernetesObject> children = new ArrayList<>();
      
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
              children.add(createProductionTestDeployment());
            }
          }
          
          String url = "Impossible to know without access to the K8s API";
          Status status = new Status(productionTestEnabled, conferenceReady, url);

          ControllerOutput desiredState = new ControllerOutput(children, status);

          log.info("> Desired State: " + desiredState);
          return desiredState;
        });
    };
  }

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

  public KubernetesObject createProductionTestDeployment() {
    
    // https://github.com/kubernetes-client/java/blob/master/fluent/src/main/java/io/kubernetes/client/openapi/models/V1DeploymentBuilder.java
    
    Map<String, String> labels = new HashMap<String, String>();
    labels.put("app", "production-tests");
    
    List<V1Container> containers = new ArrayList<V1Container>();
    containers.add(new V1Container().name("production-tests")
      .image("salaboy/metacontroller-production-tests:metacontroller")
      .imagePullPolicy("Always"));

    KubernetesObject deployment = 
      new V1Deployment()
        .apiVersion("apps/v1")
        .kind("Deployment")
        .metadata(new V1ObjectMeta()
          .name("metacontroller-production-tests")
          .labels(labels))
        .spec(new V1DeploymentSpec()
          .replicas(1)
          .selector(new V1LabelSelector()
            .matchLabels(labels))
          .template(new V1PodTemplateSpec()
            .metadata(new V1ObjectMeta().labels(labels))
            .spec(new V1PodSpec()
              .containers(containers))));
    
    return deployment;

  }
  
}
