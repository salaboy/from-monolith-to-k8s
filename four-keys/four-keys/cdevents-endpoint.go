package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-kit/log"
	"os"
	cdevents "github.com/cdevents/sdk-go/pkg/api"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	_ "github.com/lib/pq"
	"github.com/go-kit/log/level"
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

func main() {

	// The default client is HTTP.
	c, err := cloudevents.NewClientHTTP()
	if err != nil {
		level.Error(logger).Log("failed to create client, %v", err)
		return
	}
	level.Error(logger).Log(c.StartReceiver(context.Background(), receiveCDevent))
}

func receiveCDevent(event cloudevents.Event) {

	eventType, err := cdevents.ParseType(event.Type())
	if err != nil {
		level.Debug(logger).Log("this is not a valid CDEvent", eventType)
		return
	}

	// do something with event.
	//level.Debug(logger).Log("ce event", fmt.Sprintf("%s", event))

	level.Debug(logger).Log("json ce event data",  event.Data())


	cdEvent, err := cdevents.NewFromJsonBytes(event.Data())
	if err != nil {
		level.Debug(logger).Log("Failed to get CDevent from event.data() -> err: ", err)
		level.Debug(logger).Log("Failed to get CDevent from event.data()", cdEvent)
		return
	}

	level.Debug(logger).Log("cdevent ",  cdEvent)
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

	jsonCDEvent, err := json.Marshal(cdEvent)

	if err != nil {
		level.Debug(logger).Log("Failed to marshal CDevent to json ", err)
		return
	}

	// insert
	insertStmt := `insert into "cdevents_raw"("cd_source", "cd_id", "cd_timestamp", "cd_type", "cd_subject_id", "cd_subject_source", "content") values($1, $2, $3, $4, $5 , $6, $7)`
	_, e := db.Exec(insertStmt, cdEvent.GetSource(), cdEvent.GetId(), cdEvent.GetTimestamp(), cdEvent.GetType(), cdEvent.GetSubjectId(), cdEvent.GetSubjectSource(), jsonCDEvent)
	CheckError(e)

}


