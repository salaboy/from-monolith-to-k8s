package com.salaboy.controller.conference;

import io.fabric8.kubernetes.api.model.apps.Deployment;
import io.fabric8.kubernetes.client.KubernetesClient;
import io.javaoperatorsdk.operator.api.reconciler.*;
import io.javaoperatorsdk.operator.processing.dependent.kubernetes.KubernetesDependentResource;
import io.javaoperatorsdk.operator.processing.dependent.kubernetes.KubernetesDependentResourceConfig;
import io.javaoperatorsdk.operator.processing.event.source.EventSource;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.http.MediaType;
import org.springframework.web.reactive.function.client.WebClient;
import reactor.core.publisher.Mono;

import java.util.Map;
import java.util.concurrent.TimeUnit;

@ControllerConfiguration()
public class ConferenceReconciler implements Reconciler<Conference>,
        //ErrorStatusHandler<Conference>,
        EventSourceInitializer<Conference> {

    public static final String SELECTOR = "managed";

    private static final Logger log = LoggerFactory.getLogger(ConferenceReconciler.class);

    private KubernetesClient kubernetesClient;

    private WebClient.Builder webClient;

    private KubernetesDependentResource<Deployment, Conference> deploymentDR;

    public ConferenceReconciler(KubernetesClient kubernetesClient, WebClient.Builder webClient) {
        this.kubernetesClient = kubernetesClient;
        this.webClient = webClient;
        createDependentResources(kubernetesClient);
    }

    @Override
    public UpdateControl<Conference> reconcile(Conference conference, Context<Conference> context) {
        log.info("Reconciling: {}", conference.getMetadata().getName());
        conference.setStatus(new ConferenceStatus());
        String protocol = "http://";
        String servicePath = "." + conference.getSpec().getNamespace() + ".svc.cluster.local/info";
        Mono<Conference> conferenceMono = Mono.zip(getServiceInfo(protocol + "fmtok8s-frontend" + servicePath),
                        getServiceInfo(protocol + "fmtok8s-email" + servicePath),
                        getServiceInfo(protocol + "fmtok8s-agenda" + servicePath),
                        getServiceInfo(protocol + "fmtok8s-c4p" + servicePath))
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
                        if (conference.getSpec().isProductionTestsEnabled()) {
                            deploymentDR.reconcile(conference, context);
                            conference.getStatus().setProdTests(true);
                        }
                    }
                    return conference;
                });

        Conference updatedConference = conferenceMono.block();
        return UpdateControl.patchStatus(updatedConference).rescheduleAfter(5, TimeUnit.SECONDS);

    }

    private void createDependentResources(KubernetesClient client) {

        this.deploymentDR = new DeploymentDependentResource();
        deploymentDR.setKubernetesClient(client);
        deploymentDR.configureWith(new KubernetesDependentResourceConfig()
                .setLabelSelector(SELECTOR + "=true"));

    }


    @Override
    public Map<String, EventSource> prepareEventSources(EventSourceContext<Conference> context) {
        return EventSourceInitializer.nameEventSources(
                deploymentDR.initEventSource(context));
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

//    @Override
//    public ErrorStatusUpdateControl<Conference> updateErrorStatus(Conference conference, Context<Conference> context, Exception e) {
//        conference.getStatus().setErrorMessage("Error: " + e.getMessage());
//        return ErrorStatusUpdateControl.updateStatus(conference);
//    }
}
