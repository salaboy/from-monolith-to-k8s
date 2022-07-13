package com.salaboy.controller.conference;

import io.fabric8.kubernetes.model.annotation.PrinterColumn;

public class ConferenceStatus {

    @PrinterColumn(name = "READY")
    private boolean ready;
    @PrinterColumn(name = "URL")
    private String URL;
    @PrinterColumn(name = "FRONTEND")
    private boolean frontendReady;
    @PrinterColumn(name = "AGENDA")
    private boolean agendaServiceReady;
    @PrinterColumn(name = "C4P")
    private boolean c4pServiceReady;
    @PrinterColumn(name = "EMAIL")
    private boolean emailServiceReady;
    @PrinterColumn(name = "PROD TESTS")
    private boolean prodTests;

    public boolean getReady() {
        return ready;
    }

    public String getURL() {
        return URL;
    }

    public void setReady(boolean ready) {
        this.ready = ready;
    }

    public void setURL(String URL) {
        this.URL = URL;
    }

    public boolean getFrontendReady() {
        return frontendReady;
    }

    public void setFrontendReady(boolean frontendReady) {
        this.frontendReady = frontendReady;
    }

    public boolean getAgendaServiceReady() {
        return agendaServiceReady;
    }

    public void setAgendaServiceReady(boolean agendaServiceReady) {
        this.agendaServiceReady = agendaServiceReady;
    }

    public boolean getC4pServiceReady() {
        return c4pServiceReady;
    }

    public void setC4pServiceReady(boolean c4pServiceReady) {
        this.c4pServiceReady = c4pServiceReady;
    }

    public boolean getEmailServiceReady() {
        return emailServiceReady;
    }

    public void setEmailServiceReady(boolean emailServiceReady) {
        this.emailServiceReady = emailServiceReady;
    }

    public boolean getProdTests() {
        return prodTests;
    }

    public void setProdTests(boolean prodTests) {
        this.prodTests = prodTests;
    }
}
