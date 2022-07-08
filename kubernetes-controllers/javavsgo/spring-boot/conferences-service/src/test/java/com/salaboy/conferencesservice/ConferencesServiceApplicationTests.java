package com.salaboy.conferencesservice;

import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.http.MediaType;
import org.springframework.test.web.reactive.server.WebTestClient;

import java.text.ParseException;
import java.text.SimpleDateFormat;
import java.util.Arrays;
import java.util.List;

@SpringBootTest(classes = ConferencesServiceApplication.class, webEnvironment = SpringBootTest.WebEnvironment.RANDOM_PORT)
class ConferencesServiceApplicationTests {

	@Autowired
	private WebTestClient webTestClient;

	@Test
	void testGetConferences() throws ParseException {
		List<Conference> conferences = Arrays.asList(
				new Conference("123", "JBCNConf", new SimpleDateFormat("dd/MM/yy").parse("18/07/22"), "Barcelona, Spain"),
				new Conference("456", "KubeCon", new SimpleDateFormat("dd/MM/yy").parse("24/10/22"), "Detroit, USA"));


		getAll()
				.expectStatus()
				.isOk()
				.expectHeader().contentType(MediaType.APPLICATION_JSON)
				.expectBodyList(Conference.class)
				.contains(conferences.get(0), conferences.get(1))
				.returnResult()
				.getResponseBody();
	}

	private WebTestClient.ResponseSpec getAll() {
		return webTestClient.get().uri("/conferences").exchange();
	}




}
