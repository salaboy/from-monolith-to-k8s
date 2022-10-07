package domain.controlleroutput;

import com.fasterxml.jackson.annotation.JsonProperty;

public record Status(@JsonProperty("prod-tests") boolean productionTestEnabled,
                     boolean ready,
                     String url) {
}
