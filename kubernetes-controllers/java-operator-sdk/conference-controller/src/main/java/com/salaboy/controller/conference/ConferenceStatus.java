package com.salaboy.controller.conference;

public class ConferenceStatus {
    private Boolean ready;
    private String URL;

    public Boolean getReady() {
        return ready;
    }

    public String getURL() {
        return URL;
    }

    public void setReady(Boolean ready) {
        this.ready = ready;
    }

    public void setURL(String URL) {
        this.URL = URL;
    }
}
