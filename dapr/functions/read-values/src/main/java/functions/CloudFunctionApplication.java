package functions;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.annotation.Bean;
import org.springframework.messaging.Message;
import org.springframework.web.reactive.function.client.WebClient;

import java.util.function.Function;

@SpringBootApplication
public class CloudFunctionApplication {

  public static void main(String[] args) {
    SpringApplication.run(CloudFunctionApplication.class, args);
  }

  @Bean
  public Function<String, String> echo() {
    return (inputMessage) -> {
      WebClient client = WebClient.create();

      WebClient.ResponseSpec responseSpec = client.get()
          .uri("http://example.com")
          .retrieve();

      return "OK echo";
    };
  }

  @Bean
  public Function<Message<String>, String> scheduled() {
    return (inputMessage) -> {
      System.out.println("Hello from Scheduled!");
      return "OK";
    };
  }
}
