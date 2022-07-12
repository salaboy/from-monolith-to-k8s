package functions;

import java.util.Arrays;
import java.util.HashMap;
import java.util.Map;
import java.util.function.Function;

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

  @Autowired
  private WebClient.Builder webClient;
  
  @Bean
  public Function<Map<String, Object>, Mono<Map<String, Object>>> reconcile() {
    return (resource) -> {
      Map<String, Object> parent = (Map<String, Object>) resource.get("parent");
      Map<String, Object> parentMetadata = (Map<String, Object>) parent.get("metadata");
      Map<String, Object> parentSpec = (Map<String, Object>) parent.get("spec");

      log.info("Reconciling Resource: " + parent.get("apiVersion") + "/" + parent.get("Kind") + " > " + parentMetadata.get("name"));

      boolean productionTestEnabled = (boolean) parentSpec.get("production-test-enabled");

      Map<String, Object> desiredState = new HashMap<>();

      if (productionTestEnabled) {
        Map<String, Object> deployment = createProductionTestDeployment();
        desiredState.put("children", Arrays.asList(deployment));
      }

      return Mono.zip(getServiceInfo("http://fmtok8s-frontend.staging.svc.cluster.local/info"),
          getServiceInfo("http://fmtok8s-email.staging.svc.cluster.local/info"),
          getServiceInfo("http://fmtok8s-agenda.staging.svc.cluster.local/info"),
          getServiceInfo("http://fmtok8s-c4p.staging.svc.cluster.local/info"))
        .map(serviceInfos -> {
          log.info("Service Infos: " + serviceInfos);
          Map<String, Object> status = new HashMap<>();
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

          status.put("frontend-ready", frontendReady);
          status.put("email-service-ready", emailServiceReady);
          status.put("agenda-service-ready", agendaServiceReady);
          status.put("c4p-service-ready", c4pServiceReady);

          status.put("prod-tests", productionTestEnabled);

          boolean conferenceReady = false;
          if (frontendReady && emailServiceReady && agendaServiceReady && c4pServiceReady) {
            conferenceReady = true;
          }
          status.put("ready", conferenceReady);

          desiredState.put("status", status);
          status.put("url", "Impossible to know without access to the K8s API");

          log.info("> Desired State: " + desiredState);
          return desiredState;
        });
    };
  }


  
  public Mono<String> getServiceInfo(String url) {
    return webClient.build()
      .get()
      .uri(url)
      .accept(MediaType.APPLICATION_JSON)
      .retrieve()
      .bodyToMono(String.class)
      .onErrorResume(err -> Mono.just("N/A"));

  }

  public Map<String, Object> createProductionTestDeployment(){
    Yaml yaml = new Yaml();
    String deploymentYaml = "apiVersion: apps/v1\n" +
      "kind: Deployment\n" +
      "metadata:\n" +
      "  name: production-tests\n" +
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
      "          image: salaboy/production-tests:metacontroller\n" +
      "          imagePullPolicy: Always\n";
    return yaml.load(deploymentYaml);
  }
}
