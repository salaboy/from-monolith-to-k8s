package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/gorilla/mux"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	_ "github.com/lib/pq"
	"time"
)

const (
	port   = 5432
	user   = "postgres"
	dbname = "postgres"
)

type Deployment struct{

}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/deploy-frequency/day", DeploymentsByDayHandler).Methods("GET")
	log.Printf("Four Keys Metrics Server Started in 8080!")
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Returns the deployments frequency per day
func DeploymentsByDayHandler(writer http.ResponseWriter, request *http.Request) {

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

	defer rows.Close()

	for rows.Next() {
		var count int
		var time time.Time

		err = rows.Scan(&time, &count)
		CheckError(err)

		fmt.Println("Deployment: ",time, string(count))
	}


	respondWithJSON(writer, http.StatusOK, "working")
}


func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
