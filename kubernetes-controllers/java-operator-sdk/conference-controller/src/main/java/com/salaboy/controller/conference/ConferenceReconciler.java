package com.salaboy.controller.conference;

import io.fabric8.kubernetes.api.model.ListOptions;
import io.fabric8.kubernetes.api.model.ListOptionsBuilder;
import io.fabric8.kubernetes.api.model.Service;
import io.fabric8.kubernetes.api.model.ServiceList;
import io.fabric8.kubernetes.client.KubernetesClient;
import io.javaoperatorsdk.operator.api.reconciler.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.reactive.function.client.WebClient;

@ControllerConfiguration
public class ConferenceReconciler implements Reconciler<Conference> {

    private static final Logger log = LoggerFactory.getLogger(ConferenceReconciler.class);

    @Autowired
    private KubernetesClient kubernetesClient;
    @Autowired
    private WebClient.Builder webClient;

    @Override
    public UpdateControl<Conference> reconcile(Conference conference, Context context) {
        log.info("Reconciling: {}", conference.getMetadata().getName());
        //Implement logic
        ServiceList list = kubernetesClient.services().list();
        for (Service s : list.getItems()){

        }
        //webClient.build().get().uri();
        return UpdateControl.updateResource(conference);
    }

}
