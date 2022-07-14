package com.salaboy.controller.conference;

import io.fabric8.kubernetes.client.KubernetesClient;
import io.javaoperatorsdk.operator.Operator;
import io.javaoperatorsdk.operator.api.reconciler.Reconciler;
import io.javaoperatorsdk.operator.config.runtime.DefaultConfigurationService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.web.reactive.function.client.WebClient;

import java.util.List;

@Configuration
public class Config {

    @Autowired
    private KubernetesClient kubernetesClient;

    @Autowired
    private WebClient.Builder webClient;

    @Bean
    public ConferenceReconciler customServiceController(KubernetesClient kubernetesClient, WebClient.Builder webClient) {
        return new ConferenceReconciler(kubernetesClient, webClient);
    }

    @Bean(initMethod = "start", destroyMethod = "stop")
    @SuppressWarnings("rawtypes")
    public Operator operator(List<Reconciler> controllers) {
        Operator operator = new Operator();
        controllers.forEach(operator::register);
        return operator;
    }
}
