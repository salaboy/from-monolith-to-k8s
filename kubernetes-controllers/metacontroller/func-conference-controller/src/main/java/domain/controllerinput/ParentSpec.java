package domain.controllerinput;

import com.fasterxml.jackson.annotation.JsonProperty;

public record ParentSpec(String namespace, @JsonProperty("production-test-enabled") boolean productionTestEnabled) {
}
