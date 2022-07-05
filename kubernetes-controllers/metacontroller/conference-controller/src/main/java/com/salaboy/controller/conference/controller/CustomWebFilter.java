package com.salaboy.controller.conference.controller;

import org.springframework.stereotype.Component;
import org.springframework.web.server.ServerWebExchange;
import org.springframework.web.server.WebFilter;
import org.springframework.web.server.WebFilterChain;
import reactor.core.publisher.Mono;

@Component
public class CustomWebFilter implements WebFilter {
    @Override
    public Mono<Void> filter(ServerWebExchange serverWebExchange,
                             WebFilterChain webFilterChain) {
        BodyCaptureExchange bodyCaptureExchange = new BodyCaptureExchange(serverWebExchange);
        return webFilterChain.filter(bodyCaptureExchange).doOnSuccess( (se) -> {
            System.out.println("Headers request "+bodyCaptureExchange.getRequest().getHeaders());
            System.out.println("Body request "+bodyCaptureExchange.getRequest().getFullBody());
            System.out.println("Headers response "+bodyCaptureExchange.getResponse().getHeaders());
            System.out.println("Body response "+bodyCaptureExchange.getResponse().getFullBody());
        });
    }
}
