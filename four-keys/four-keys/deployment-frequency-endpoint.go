package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"time"
)

const (
	port   = 5432
	user   = "postgres"
	dbname = "postgres"
)

type DeploymentFrequency struct{
	DeployName string
	Deployments int
	Time time.Time
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

	rows, err := db.Query(`SELECT distinct deploy_name AS NAME, DATE_TRUNC('day', time_created) AS day, 
						COUNT(distinct deploy_id) AS deployments FROM deployments GROUP BY deploy_name, day`)
	CheckError(err)

	defer rows.Close()
	deployments  := make([]DeploymentFrequency, 0)
	for rows.Next() {
		var deployName string
		var count int
		var time time.Time

		err = rows.Scan(&deployName, &time, &count)
		CheckError(err)

	    deploy :=  DeploymentFrequency{
			DeployName:  deployName,
			Deployments: count,
			Time:        time,
		}

	    deployments = append(deployments, deploy)
	}

	respondWithJSON(writer, http.StatusOK, deployments)
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
