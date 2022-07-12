package com.salaboy.controller.conference;

public class ConferenceStatus {
    private Boolean ready;
    private String URL;
    private Boolean frontendReady;
    private Boolean agendaServiceReady;
    private Boolean c4pServiceReady;
    private Boolean emailServiceReady;
    private Boolean prodTests;

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

    public Boolean getFrontendReady() {
        return frontendReady;
    }

    public void setFrontendReady(Boolean frontendReady) {
        this.frontendReady = frontendReady;
    }

    public Boolean getAgendaServiceReady() {
        return agendaServiceReady;
    }

    public void setAgendaServiceReady(Boolean agendaServiceReady) {
        this.agendaServiceReady = agendaServiceReady;
    }

    public Boolean getC4pServiceReady() {
        return c4pServiceReady;
    }

    public void setC4pServiceReady(Boolean c4pServiceReady) {
        this.c4pServiceReady = c4pServiceReady;
    }

    public Boolean getEmailServiceReady() {
        return emailServiceReady;
    }

    public void setEmailServiceReady(Boolean emailServiceReady) {
        this.emailServiceReady = emailServiceReady;
    }

    public Boolean getProdTests() {
        return prodTests;
    }

    public void setProdTests(Boolean prodTests) {
        this.prodTests = prodTests;
    }
}
