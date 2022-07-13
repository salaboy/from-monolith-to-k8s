package com.salaboy;

import io.quarkus.test.junit.QuarkusTest;
import org.junit.jupiter.api.Test;

import javax.ws.rs.core.MediaType;
import java.text.ParseException;
import java.text.SimpleDateFormat;
import java.util.Arrays;
import java.util.Date;
import java.util.List;

import static io.restassured.RestAssured.given;
import static org.hamcrest.CoreMatchers.is;
import static org.hamcrest.Matchers.*;

@QuarkusTest
public class ConferenceResourceTest {

    @Test
    public void testConferenceEndpoint() {

        given()
          .when().get("/conferences")
          .then()
             .statusCode(200)
                .contentType(MediaType.APPLICATION_JSON)
                .body("$.size()", is(2),
                        "[0].id", is("123"),
                        "[0].name", is("JBCNConf"),
                        "[0].where", is("Barcelona, Spain"),
                        "[1].id", is("456"),
                        "[1].name", is("KubeCon"),
                        "[1].where", is("Detroit, USA")
                        );
    }

}