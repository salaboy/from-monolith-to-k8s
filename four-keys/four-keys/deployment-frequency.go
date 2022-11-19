package main

import (
	"database/sql"
	"encoding/json"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	cdevents "github.com/cdevents/sdk-go/pkg/api"

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

// reads from CDEvents and look for Deployments to count them
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
	// @TODO: add a way to keep track of which ones were processed
	rows, err := db.Query(`SELECT "cd_id", "content" FROM "cdevents_raw" where cd_type=$1`,  cdevents.ServiceDeployedEventV1.String())
	CheckError(err)

	defer rows.Close()
	for rows.Next() {
		var id string
		var content []byte

		err = rows.Scan(&id, &content)
		CheckError(err)


		event, _ := cdevents.NewServiceDeployedEvent()

		err := json.Unmarshal(content, &event)
		if err != nil {
			level.Debug(logger).Log("Ignoring id ", id)
			continue
		}

		// insert
		insertStmt := `insert into "deployments"("deploy_id", "deploy_name", "time_created") values($1, $2, $3 )`
		_, e := db.Exec(insertStmt, event.Subject.Content.ArtifactId, event.GetSubjectId(), event.GetTimestamp())

		if e != nil {
			level.Error(logger).Log("Inserting failed for event ", e)
		}

	}

	CheckError(err)

}

