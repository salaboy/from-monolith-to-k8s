package com.salaboy.dapr.javaapp;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.RestController;

import io.dapr.client.DaprClient;
import io.dapr.client.DaprClientBuilder;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import io.dapr.client.domain.State;


import java.util.List;
import java.util.ArrayList;

@SpringBootApplication
@RestController
public class JavaAppApplication {

	private static final Logger log = LoggerFactory.getLogger(JavaAppApplication.class);

	private static final String STATE_STORE_NAME = "my-db-dapr-statestore";

	private DaprClient client = new DaprClientBuilder().build();

	public static void main(String[] args) {
		SpringApplication.run(JavaAppApplication.class, args);

	}

	@PostMapping("/")
	public MyValues storeValues(@RequestParam("value") String value) {
		State<MyValues> results = client.getState(STATE_STORE_NAME, "values", MyValues.class).block();

		MyValues valuesList = results.getValue();

		if (valuesList == null) {
			valuesList = new MyValues(new ArrayList<String>());
			valuesList.values().add(value);
		} else {
			valuesList.values().add(value);
		}
	
		client.saveState(STATE_STORE_NAME, "values", valuesList).block();
		return valuesList;
	}

	@DeleteMapping("/")
	public void deleteAllValues() {
		client.deleteState(STATE_STORE_NAME, "values").block();
	}

	public record MyValues(List<String> values) {}

}


