package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/go-kit/log/level"
	"github.com/go-kit/log"

	cdevents "github.com/cdevents/sdk-go/pkg/api"
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

func CheckError(err error) {
	if err != nil {
		level.Error(logger).Log("error", err)
		panic(err)
	}
}

var client cloudevents.Client


func main() {

	client, err := cloudevents.NewClientHTTP()
	if err != nil {
		level.Error(logger).Log("failed to create client, %v", err)
		return
	}

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


		if event.Type() == cdevents.ArtifactPackagedEventV1.String() {
			// send request to mapped function(s)
		}else if event.Type() == "dev.knative.apiserver.resource.add" {

			// @TODO: get URL from cloudevent type from the routing table
			var functionUrl = "http://api-server-to-service-deployment-function.four-keys.svc.cluster.local"

			level.Debug(logger).Log("Sending event to function: ", functionUrl)
			// send request to mapped function(s)
			ctx := cloudevents.ContextWithTarget(context.Background(), functionUrl)
			ctx = cloudevents.WithEncodingBinary(ctx)
			if result := client.Send(ctx, event); cloudevents.IsUndelivered(result) {
				level.Debug(logger).Log("failed to send:", result)
			}
		}



	}

	CheckError(err)

}

