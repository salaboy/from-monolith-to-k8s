package com.salaboy.controller.conference.controller;


import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.MediaType;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.reactive.function.client.WebClient;
import org.yaml.snakeyaml.Yaml;
import reactor.core.publisher.Mono;

import java.util.HashMap;
import java.util.Map;

@RestController
@RequestMapping
public class ConferenceController {

    private static final Logger log = LoggerFactory.getLogger(ConferenceController.class);

    @Autowired
    private WebClient.Builder webClient;

    @PostMapping(produces = MediaType.APPLICATION_JSON_VALUE)
    public Mono<Map<String, Object>> reconcileResource(@RequestBody Map<String, Object> resource) throws JsonProcessingException {
        log.info("> REST ENDPOINT INVOKED for reconciling Resource: " + resource);
        log.info("> Resource Parent: " + resource.get("parent"));
        log.info("> Resource Children: " + resource.get("children"));

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
        Map<String, Object> deployment = yaml.load(deploymentYaml);

        Map<String, Object> desiredState = new HashMap<>();

        Map<String, Object> status = new HashMap<>();



        status.put("ready", true);
        status.put("url", "1.2.3.4");
        desiredState.put("status", status);
        desiredState.put("children", deployment);
        return Mono.just(desiredState);
    }
}
