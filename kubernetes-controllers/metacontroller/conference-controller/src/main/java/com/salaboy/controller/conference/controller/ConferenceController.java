package com.salaboy.controller.conference.controller;


import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.http.MediaType;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import reactor.core.publisher.Mono;

import java.util.HashMap;
import java.util.Map;

@RestController
@RequestMapping
public class ConferenceController {

    private static final Logger log = LoggerFactory.getLogger(ConferenceController.class);

    @PostMapping(produces = MediaType.APPLICATION_JSON_VALUE)
    public Mono<Map<String, Object>> reconcileResource(@RequestBody Map<String, Object> resource) {
        log.info("> REST ENDPOINT INVOKED for reconciling Resource: " + resource);
        Map<String, Object> desiredState = new HashMap<>();
        Map<String, Object> status = new HashMap<>();
        status.put("ready", true);
        status.put("url", "1.2.3.4");
        desiredState.put("status", status);
        return Mono.just(desiredState);
    }
}
