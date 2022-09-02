package main

import (
	"database/sql"
	"encoding/json"
	cdevents "github.com/cdevents/sdk-go/pkg/api"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"log"

	"fmt"
	_ "github.com/lib/pq"

	"os"
)

const (
	port   = 5432
	user   = "postgres"
	dbname = "postgres"
)

// cdevents_raw table `CREATE TABLE IF NOT EXISTS cdevents_raw ( cd_source varchar(255) NOT NULL, cd_id varchar(255) NOT NULL, cd_type varchar(255) NOT NULL, cd_subject_id varchar(255) NOT NULL,cd_subject_type varchar(255), cd_subject_source varchar(255), content json NOT NULL, PRIMARY KEY (cd_source, cd_id));`
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

	fmt.Println("Connected!")

	// select
	rows, err := db.Query(`SELECT "event_id", "content" FROM "cloudevents_raw"`) // add a way to keep track of which ones were processed
	CheckError(err)

	defer rows.Close()
	for rows.Next() {
		var id int
		var content []byte

		err = rows.Scan(&id, &content)
		CheckError(err)

		fmt.Println(id, string(content))

		event := cloudevents.NewEvent()

		err := json.Unmarshal(content, &event)
		if err != nil {
			fmt.Println("Ignoring id ", id)
			continue
		}

		//Maybe check if it is a CDEvent first, before trying to map it
		// cdevents.IsACDevent(event) -> true then insert

		cdEvent, err := mapToCDEvents(event)
		if err != nil {
			fmt.Println("There is no mapping for this type of event", event.Type(), " - Ignoring id ", id)
			continue
		}


		// insert
		insertStmt := `insert into "cdevents_raw"("cd_source", "cd_id", "cd_type", "cd_subject_id", "content") values($1, $2, $3, $4, $5 )`
		_, e := db.Exec(insertStmt, cdEvent.GetSource(), cdEvent.GetId(), cdEvent.GetType(), cdEvent.GetSubjectId(), content)

		if e != nil {
			fmt.Println("Inserting failed for event ", e)
		}

	}

	CheckError(err)

}

// This can be splitted into functions, so one function Mapping
///   Where mapping means CloudEvent -> CDEvent,
//      this can be done by looking for certain Cloud Event types and creating a CDEvent or even by parsing the CloudEvent boyd and selecting which proerties to use to create a CDEvent
func mapToCDEvents(event cloudevents.Event) (cdevents.CDEvent, error) {
    // Define mappings for use case

	// CloudEvents to CDEvents Mapping goes here
	if event.Type() == cdevents.ArtifactPackagedEventV1.String() {
		event, err := cdevents.NewPipelineRunQueuedEvent()
		if err != nil {
			log.Fatalf("could not create a cdevent, %v", err)
			return nil, err
		}
		// Set the required context fields
		event.SetSubjectId("myPipelineRun1")
		event.SetSource("my/first/cdevent/program")

		// Set the required subject fields
		event.SetSubjectPipelineName("myPipeline")
		event.SetSubjectUrl("https://example.com/myPipeline")
		return event, nil
	}
	return nil, nil
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
