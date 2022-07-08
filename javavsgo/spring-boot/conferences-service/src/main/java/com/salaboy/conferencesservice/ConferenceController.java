package com.salaboy.conferencesservice;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import reactor.core.publisher.Flux;

import java.text.ParseException;
import java.text.SimpleDateFormat;
import java.util.Date;

@RestController()
@RequestMapping("/conferences")
public class ConferenceController {

    @GetMapping
    public Flux<Conference> getAllConference() throws ParseException {
        return Flux.just(new Conference("123","JBCNConf", new SimpleDateFormat("dd/MM/yy").parse("18/07/22"), "Barcelona, Spain" ),
                new Conference("456","KubeCon", new SimpleDateFormat("dd/MM/yy").parse("24/10/22"), "Detroit, USA" ));
    }
}

record Conference(String id, String name, Date when, String where){}