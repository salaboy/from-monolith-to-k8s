package domain.controllerinput;

public record Parent(String apiVersion, String kind, ParentMetadata metadata, ParentSpec spec) {
}
