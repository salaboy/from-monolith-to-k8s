package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"os"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	_ "github.com/lib/pq"
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

// cloudevents_raw table `CREATE TABLE IF NOT EXISTS cloudevents_raw ( event_id serial NOT NULL PRIMARY KEY, content json NOT NULL, event_timestamp TIMESTAMP NOT NULL);`
func main() {

	// The default client is HTTP.
	c, err := cloudevents.NewClientHTTP()
	if err != nil {
		level.Error(logger).Log("failed to create client, %v", err)
	}
	level.Error(logger).Log(c.StartReceiver(context.Background(), receiveCloudEvent))
}

func receiveCloudEvent(event cloudevents.Event) {

	// do something with event.
	level.Debug(logger).Log("event", event)
	jsonEvent, err := json.Marshal(event)
	CheckError(err)

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

	// insert
	insertStmt := `insert into "cloudevents_raw"("content", "event_timestamp") values($1, current_timestamp)`
	_, e := db.Exec(insertStmt, string(jsonEvent))
	if e != nil {
		level.Error(logger).Log("Inserting failed for event ", e)
	}
	CheckError(e)

}


