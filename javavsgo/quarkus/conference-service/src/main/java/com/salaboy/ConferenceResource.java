package com.salaboy;

import javax.ws.rs.GET;
import javax.ws.rs.Path;
import javax.ws.rs.Produces;
import javax.ws.rs.core.MediaType;
import java.text.ParseException;
import java.text.SimpleDateFormat;
import java.util.Arrays;
import java.util.List;

@Path("/conferences")
public class ConferenceResource {

    @GET
    @Produces(MediaType.APPLICATION_JSON)
    public List<Conference> getAllConferences() throws ParseException {
        return Arrays.asList(
                new Conference("123","JBCNConf", new SimpleDateFormat("dd/MM/yy").parse("18/07/22"), "Barcelona, Spain" ),
                new Conference("456","KubeCon",  new SimpleDateFormat("dd/MM/yy").parse("24/10/22"), "Detroit, USA" ));
    }
}
