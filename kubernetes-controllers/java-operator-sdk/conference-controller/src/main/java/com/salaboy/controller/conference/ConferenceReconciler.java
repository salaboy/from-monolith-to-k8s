package com.salaboy.controller.conference;

import io.fabric8.kubernetes.api.model.apps.Deployment;
import io.fabric8.kubernetes.client.KubernetesClient;
import io.javaoperatorsdk.operator.api.reconciler.Context;
import io.javaoperatorsdk.operator.api.reconciler.ControllerConfiguration;
import io.javaoperatorsdk.operator.api.reconciler.Reconciler;
import io.javaoperatorsdk.operator.api.reconciler.UpdateControl;
import io.javaoperatorsdk.operator.api.reconciler.dependent.Dependent;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.MediaType;
import org.springframework.web.reactive.function.client.WebClient;
import reactor.core.publisher.Mono;

@ControllerConfiguration(dependents = {
        @Dependent(type = DeploymentDependentResource.class),
    })
public class ConferenceReconciler implements Reconciler<Conference> {

    public static final String SELECTOR = "managed";

    private static final Logger log = LoggerFactory.getLogger(ConferenceReconciler.class);

    @Autowired
    private KubernetesClient kubernetesClient;

    @Autowired
    private WebClient.Builder webClient;

    @Override
    public UpdateControl<Conference> reconcile(Conference conference, Context<Conference> context) {
        log.info("Reconciling: {}", conference.getMetadata().getName());

        conference.setStatus(new ConferenceStatus());

        if (conference.getSpec().isProductionTestsEnabled()) {
            final var name = context.getSecondaryResource(Deployment.class).orElseThrow().getMetadata().getName();
            log.info("Deployment Created with name: " + name);
            conference.getStatus().setProdTests(true);
        }

        Mono<Conference> conferenceMono = Mono.zip(getServiceInfo("http://fmtok8s-frontend.staging.svc.cluster.local/info"),
                        getServiceInfo("http://fmtok8s-email.staging.svc.cluster.local/info"),
                        getServiceInfo("http://fmtok8s-agenda.staging.svc.cluster.local/info"),
                        getServiceInfo("http://fmtok8s-c4p.staging.svc.cluster.local/info"))
                .map(serviceInfos -> {
                    log.info("Service Infos: " + serviceInfos);


                    if (!serviceInfos.getT1().contains("N/A") && !serviceInfos.getT1().isEmpty()) {
                        conference.getStatus().setFrontendReady(true);
                    }
                    if (!serviceInfos.getT2().contains("N/A") && !serviceInfos.getT2().isEmpty()) {
                        conference.getStatus().setEmailServiceReady(true);
                    }
                    if (!serviceInfos.getT3().contains("N/A") && !serviceInfos.getT3().isEmpty()) {
                        conference.getStatus().setAgendaServiceReady(true);
                    }
                    if (!serviceInfos.getT4().contains("N/A") && !serviceInfos.getT4().isEmpty()) {
                        conference.getStatus().setC4pServiceReady(true);
                    }

                    if (conference.getStatus().getFrontendReady() &&
                            conference.getStatus().getEmailServiceReady() &&
                            conference.getStatus().getAgendaServiceReady() &&
                            conference.getStatus().getC4pServiceReady()) {
                        conference.getStatus().setReady(true);
                    }

                    return conference;
                });

        Conference updatedConference = conferenceMono.block();
        return UpdateControl.patchStatus(updatedConference);

    }

//    private void createProductionTestDeployment() {
//        Deployment deployment = new DeploymentBuilder()
//                .withNewMetadata()
//                .withName("production-tests")
//                .endMetadata()
//                .withNewSpec()
//                .withReplicas(1)
//                .withNewTemplate()
//                .withNewMetadata()
//                .addToLabels("app", "production-tests")
//                .endMetadata()
//                .withNewSpec()
//                .addNewContainer()
//                .withName("production-tests")
//                .withImage("salaboy/production-tests:java-operator-sdk")
//                .withImagePullPolicy("Always")
//                .addNewPort()
//                .withContainerPort(8080)
//                .endPort()
//                .endContainer()
//                .endSpec()
//                .endTemplate()
//                .withNewSelector()
//                .addToMatchLabels("app", "production-tests")
//                .endSelector()
//                .endSpec()
//                .build();
//
//        deployment = kubernetesClient.apps().deployments().inNamespace("default").create(deployment);
//        log.info("Created deployment: {}", deployment);
//    }

    public Mono<String> getServiceInfo(String url) {
        return webClient.build()
                .get()
                .uri(url)
                .accept(MediaType.APPLICATION_JSON)
                .retrieve()
                .bodyToMono(String.class)
                .onErrorResume(err -> Mono.just("N/A"));

    }
}
