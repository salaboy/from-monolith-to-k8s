package main

import (
    "context"
    "database/sql"
    "encoding/json"
    "fmt"
    _ "github.com/lib/pq"
    cloudevents "github.com/cloudevents/sdk-go/v2"
    "log"
    "os"
)
 
const (
    port     = 5432
    user     = "postgres"
    dbname   = "postgres"
)



//  cloudevents_raw table `CREATE TABLE IF NOT EXISTS cloudevents_raw ( event_id serial NOT NULL PRIMARY KEY, content json NOT NULL, event_timestamp TIMESTAMP NOT NULL);`
func main() {

    // The default client is HTTP.
    c, err := cloudevents.NewClientHTTP()
    if err != nil {
        log.Fatalf("failed to create client, %v", err)
    }
    log.Fatal(c.StartReceiver(context.Background(), receive))
}

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


    // insert
    insertStmt := `insert into "cloudevents_raw"("content", "event_timestamp") values($1, current_timestamp)`
    _, e := db.Exec(insertStmt, string(jsonEvent))
    CheckError(e)

}
 
func CheckError(err error) {
    if err != nil {
        panic(err)
    }
}
