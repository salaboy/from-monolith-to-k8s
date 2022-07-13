package com.salaboy.controller.conference;

import com.fasterxml.jackson.annotation.JsonProperty;

public class ConferenceSpec {
    @JsonProperty("production-tests-enabled")
    private Boolean productionTestsEnabled;

    public Boolean isProductionTestsEnabled() {
        return productionTestsEnabled;
    }

    public void setProductionTestsEnabled(Boolean productionTestsEnabled) {
        this.productionTestsEnabled = productionTestsEnabled;
    }
}
