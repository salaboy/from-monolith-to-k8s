package com.salaboy.controller.conference;

import io.javaoperatorsdk.operator.api.reconciler.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

@ControllerConfiguration
public class ConferenceReconciler implements Reconciler<Conference> {

    private static final Logger log = LoggerFactory.getLogger(ConferenceReconciler.class);

    @Override
    public UpdateControl<Conference> reconcile(Conference conference, Context context) {
        log.info("Reconciling: {}", conference.getMetadata().getName());
        //Implement logic
        return UpdateControl.updateResource(conference);
    }

}
