package functions;

import com.metacontrollerjava.api.input.ControllerInput;
import com.metacontrollerjava.api.output.ControllerOutput;
import com.metacontrollerjava.api.output.Status;
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
import java.util.Arrays;
import java.util.HashMap;
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

	private static boolean isServiceReady(String input) {
		return !input.contains("N/A") && !input.isEmpty();
	}

	@Bean
	public Function<ControllerInput, Mono<ControllerOutput>> reconcile(WebClient http) {
		// Input: A controller JSON object that contains a CompositeController,
		// a Parent and Child, where the Parent describes the user input that triggered
		// the function call
		// - The CompositeController is controller responsible for the webhook
		// - The parent is the triggering input - the resource that needs to be reconciled
		// - The children are the resources that can be modified as part of the
		// reconciliation loop
		// See CompositeController.Spec for a summary of the triggering resource, the
		// hook, and the children (the resources that can be changed)
		return (resource) -> {
			var parent = resource.parent();
			var parentMetadata = parent.metadata();
			var parentSpec = parent.spec();
			log.info("Reconciling Resource: " + parent.apiVersion() + "/" + parent.kind() + " > "
					+ parentMetadata.name());
			var productionTestEnabled = parentSpec.productionTestEnabled();
			var conferenceNamespace = parentSpec.namespace();
			var children = new ArrayList<KubernetesObject>();
			return Mono
					.zip(getServiceInfo(http, "fmtok8s-frontend", conferenceNamespace),
							getServiceInfo(http, "fmtok8s-email", conferenceNamespace),
							getServiceInfo(http, "fmtok8s-agenda", conferenceNamespace),
							getServiceInfo(http, "fmtok8s-c4p", conferenceNamespace)) //
					.map(serviceInfos -> {
						log.info("Service Infos: " + serviceInfos);
						var frontendReady = isServiceReady(serviceInfos.getT1());
						var agendaServiceReady = isServiceReady(serviceInfos.getT2());
						var emailServiceReady = isServiceReady(serviceInfos.getT3());
						var c4pServiceReady = isServiceReady(serviceInfos.getT4());
						var conferenceReady = false;
						if (frontendReady && emailServiceReady && agendaServiceReady && c4pServiceReady) {
							conferenceReady = true;
							if (productionTestEnabled) {
								children.add(createProductionTestDeployment(parentMetadata.name()));
							}
						}
						var url = "Impossible to know without access to the K8s API";
						var status = new Status(productionTestEnabled, conferenceReady, url);
						var desiredState = new ControllerOutput(children, status);
						log.info("> Desired State: " + desiredState);
						return desiredState;
					});
		};
	}

	public Mono<String> getServiceInfo(WebClient http, String name, String namespace) {
		return getServiceInfo(http, "http://" + name + "." + namespace + ".svc.cluster.local/info");
	}

	public Mono<String> getServiceInfo(WebClient http, String url) {
		return http.get().uri(url).accept(MediaType.APPLICATION_JSON).retrieve().bodyToMono(String.class)
				.onErrorResume(err -> Mono.just("N/A"));

	}

	public KubernetesObject createProductionTestDeployment(String name) {

		// https://github.com/kubernetes-client/java/blob/master/fluent/src/main/java/io/kubernetes/client/openapi/models/V1DeploymentBuilder.java

		var labels = new HashMap<String, String>();
		labels.put("app", "production-tests-" + name);

		var containers = new ArrayList<V1Container>();
		containers.add(new V1Container().name("production-tests").image("alpine")
				.command(Arrays.asList(new String[] { "sh" }))
				.args(Arrays.asList(new String[] { "-c",
						"while true; do echo \"Running production tests @ \" `date`; sleep 10; done" }))
				.imagePullPolicy("Always"));

	// @formatter:off
    var deploymentKubernetesObject =
      new V1Deployment()
        .apiVersion("apps/v1")
        .kind("Deployment")
        .metadata(new V1ObjectMeta()
          .name("metacontroller-production-tests-" + name)
          .labels(labels))
        .spec(new V1DeploymentSpec()
          .replicas(1)
          .selector(new V1LabelSelector()
            .matchLabels(labels))
          .template(new V1PodTemplateSpec()
            .metadata(new V1ObjectMeta().labels(labels))
            .spec(new V1PodSpec()
              .containers(containers))));
    // @formatter:on

		return deploymentKubernetesObject;

	}

}
