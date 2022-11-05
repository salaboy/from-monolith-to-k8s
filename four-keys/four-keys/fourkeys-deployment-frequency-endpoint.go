package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	_ "github.com/lib/pq"
)

const (
	port   = 5432
	user   = "postgres"
	dbname = "postgres"
)

type Deployment struct{

}

func main() {

	// The default client is HTTP.
	c, err := cloudevents.NewClientHTTP()
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}
	log.Fatal(c.StartReceiver(context.Background(), receive))
}

// Normal Endpoint to return deployment frequency per service

func receive(event cloudevents.Event) {

	// do something with event.
	fmt.Printf("%s", event)
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

	fmt.Println("Connected!")

	rows, err := db.Query(`SELECT DATE_TRUNC('day', time_created) AS day, COUNT(distinct deploy_id) AS deployments FROM deployments GROUP BY day;`)

	CheckError(err)

}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
