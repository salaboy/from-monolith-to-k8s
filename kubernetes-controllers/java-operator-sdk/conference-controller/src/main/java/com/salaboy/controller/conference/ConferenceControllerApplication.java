package com.salaboy.controller.conference;

import io.javaoperatorsdk.operator.Operator;
import io.javaoperatorsdk.operator.config.runtime.DefaultConfigurationService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.CommandLineRunner;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication
public class ConferenceControllerApplication {

	public static void main(String[] args) {
		SpringApplication.run(ConferenceControllerApplication.class, args);
	}

}
