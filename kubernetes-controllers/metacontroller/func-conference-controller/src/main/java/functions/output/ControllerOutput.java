package functions.output;

import io.kubernetes.client.common.KubernetesObject;

import java.util.List;

public record ControllerOutput(List<KubernetesObject> children, Status status) {
}
