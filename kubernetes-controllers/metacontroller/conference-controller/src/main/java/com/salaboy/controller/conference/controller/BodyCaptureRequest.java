package com.salaboy.controller.conference.controller;

import org.springframework.core.io.buffer.DataBuffer;
import org.springframework.http.server.reactive.ServerHttpRequest;
import org.springframework.http.server.reactive.ServerHttpRequestDecorator;
import reactor.core.publisher.Flux;

import java.nio.charset.StandardCharsets;

public class BodyCaptureRequest extends ServerHttpRequestDecorator {
    private final StringBuilder body = new StringBuilder();

    public BodyCaptureRequest(ServerHttpRequest delegate) {
        super(delegate);
    }

    public Flux<DataBuffer> getBody() {
        return super.getBody().doOnNext(this::capture);
    }

    private void capture(DataBuffer buffer) {
        this.body.append(StandardCharsets.UTF_8.decode(buffer.asByteBuffer()).toString());
    }

    public String getFullBody() {
        return this.body.toString();
    }
}
