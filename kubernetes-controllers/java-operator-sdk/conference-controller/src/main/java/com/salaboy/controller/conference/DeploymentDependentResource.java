package com.salaboy.controller.conference;

import io.fabric8.kubernetes.api.model.apps.Deployment;
import io.fabric8.kubernetes.api.model.apps.DeploymentBuilder;
import io.javaoperatorsdk.operator.api.reconciler.Context;
import io.javaoperatorsdk.operator.processing.dependent.kubernetes.CRUKubernetesDependentResource;
import io.javaoperatorsdk.operator.processing.dependent.kubernetes.KubernetesDependent;

import java.util.Map;

import static com.salaboy.controller.conference.ConferenceReconciler.SELECTOR;

@KubernetesDependent(labelSelector = SELECTOR)
public class DeploymentDependentResource extends CRUKubernetesDependentResource<Deployment, Conference> {

    public DeploymentDependentResource() {
        super(Deployment.class);
    }

    @Override
    protected Deployment desired(Conference conference, Context<Conference> context) {
        Deployment deployment = new DeploymentBuilder()
                .withNewMetadata()
                .withName("java-operator-sdk-production-tests")
                .withNamespace(conference.getMetadata().getNamespace())
                .endMetadata()
                .withNewSpec()
                .withReplicas(1)
                .withNewTemplate()
                .withNewMetadata()
                .withLabels(Map.of(SELECTOR,"true"))
                .addToLabels("app", "production-tests")
                .endMetadata()
                .withNewSpec()
                .addNewContainer()
                .withName("production-tests")
                .withImage("salaboy/java-operator-sdk-production-tests:java-operator-sdk")
                .withImagePullPolicy("Always")
                .addNewPort()
                .withContainerPort(8080)
                .endPort()
                .endContainer()
                .endSpec()
                .endTemplate()
                .withNewSelector()
                .addToMatchLabels("app", "production-tests")
                .endSelector()
                .endSpec()
                .build();
        return deployment;
    }
}
