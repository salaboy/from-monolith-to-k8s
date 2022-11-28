package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/go-kit/log/level"
	"github.com/go-kit/log"
	"gopkg.in/yaml.v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"

	"fmt"

	_ "github.com/lib/pq"

	"os"
)

const (
	port   = 5432
	user   = "postgres"
	dbname = "postgres"
)

var logger = log.NewLogfmtLogger(os.Stdout)

type RoutesConfig struct{
	Routes map[string][]string `yaml:"routes"`
}


/*
    routes:
      "dev.knative.apiserver.resource.add":
      - http://api-server-to-service-deployment-function.four-keys.svc.cluster.local
      - http://sockeye.default.svc.cluster.local
      "dev.knative.github.pr.new":
      - http://sockeye.default.svc.cluster.local
 */

func CheckError(err error) {
	if err != nil {
		level.Error(logger).Log("error", err)
		panic(err)
	}
}

var client cloudevents.Client
var CE_ROUTES = os.Getenv("CE_ROUTES")

func main() {

	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", os.Getenv("POSTGRESQL_HOST"), port, user, os.Getenv("POSTGRESQL_PASS"), dbname)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	// close database
	defer db.Close()

	// check db
	err = db.Ping()
	CheckError(err)

	level.Info(logger).Log("DB", "Connected!")

	// select
	rows, err := db.Query(`SELECT "event_id", "content" FROM "cloudevents_raw"`) // add a way to keep track of which ones were processed
	CheckError(err)

	defer rows.Close()
	for rows.Next() {
		var id int
		var content []byte

		err = rows.Scan(&id, &content)
		CheckError(err)

		//level.Debug(logger).Log(id, string(content))

		event := cloudevents.NewEvent()

		err := json.Unmarshal(content, &event)
		if err != nil {
			level.Debug(logger).Log("Ignoring id ", id)

			continue
		}


		routed := routeCloudEvent(event)
		if routed == 0 {
			level.Debug(logger).Log("Event ignored, there is no routing rule defined", event.Type())
		}else{
			level.Debug(logger).Log("Event routed to functions: ", routed)
		}

	}

	CheckError(err)

}

// routeCloudEvent based on different parameters that users can set using the environment variable CE_ROUTES
func routeCloudEvent(event event.Event) int {
	//Route by Type, but other routing mechanisms can be implemented
	client, err := cloudevents.NewClientHTTP()
	if err != nil {
		level.Error(logger).Log("failed to create client, %v", err)
		return 0
	}
	routesDefinition, err := ReadRoutesConfigFromEnvString(CE_ROUTES)
	CheckError(err)

	routeList := routesDefinition.Routes[event.Type()]

	level.Debug(logger).Log("routing config for event type", fmt.Sprintf("%+v\n",routeList))

	eventsRouted := 0
	for  _, url := range routeList{
		level.Debug(logger).Log("Sending event to function: ", url)
		// send request to mapped function(s)
		ctx := cloudevents.ContextWithTarget(context.Background(), url)
		ctx = cloudevents.WithEncodingBinary(ctx)
		if result := client.Send(ctx, event); cloudevents.IsUndelivered(result) {
			level.Error(logger).Log("failed to send:", result)
		}
		eventsRouted++
	}
	return eventsRouted


}

func ReadRoutesConfigFromEnvString(routesContent string) (RoutesConfig, error){
	routesDefinition := RoutesConfig{}
	level.Debug(logger).Log("about to parse the following yaml: %s \n", routesContent)
	err := yaml.Unmarshal([]byte(routesContent), &routesDefinition)
	if err != nil {
		level.Error(logger).Log("failed to parse routes:", err)

		return RoutesConfig{}, err
	}

	level.Debug(logger).Log("Routes Definitions: ", fmt.Sprintf("%+v\n", routesDefinition))


	return routesDefinition, nil
}

