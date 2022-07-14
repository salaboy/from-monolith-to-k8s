package com.salaboy.controller.conference;

import com.fasterxml.jackson.annotation.JsonProperty;

public class ConferenceSpec {
    @JsonProperty("production-tests-enabled")
    private boolean productionTestsEnabled;
    private String namespace;

    public boolean isProductionTestsEnabled() {
        return productionTestsEnabled;
    }

    public void setProductionTestsEnabled(boolean productionTestsEnabled) {
        this.productionTestsEnabled = productionTestsEnabled;
    }

    public boolean getProductionTestsEnabled() {
        return productionTestsEnabled;
    }

    public String getNamespace() {
        return namespace;
    }

    public void setNamespace(String namespace) {
        this.namespace = namespace;
    }
}
