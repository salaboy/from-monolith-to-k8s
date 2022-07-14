package com.salaboy.controller.conference;


import io.fabric8.kubernetes.api.model.Namespaced;
import io.fabric8.kubernetes.client.CustomResource;
import io.fabric8.kubernetes.model.annotation.Group;
import io.fabric8.kubernetes.model.annotation.ShortNames;
import io.fabric8.kubernetes.model.annotation.Version;

@Group("java-operator-sdk.conference.salaboy.com")
@Version("v1")
@ShortNames("conf")
public class Conference extends CustomResource<ConferenceSpec, ConferenceStatus> implements
        Namespaced {

}
