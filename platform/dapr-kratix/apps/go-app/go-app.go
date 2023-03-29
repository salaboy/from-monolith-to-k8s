package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/gorilla/mux"
)

var (
	STATE_STORE_NAME = ""
	daprClient       dapr.Client
)

type MyValues struct {
	Values []string
}

func readValues(w http.ResponseWriter, r *http.Request) {

	STATE_STORE_NAME := os.Getenv("STATE_STORE_NAME")
	if STATE_STORE_NAME == "" {
		STATE_STORE_NAME = "my-dapr-db-statestore"
	}

	ctx := context.Background()

	daprClient, daprErr := dapr.NewClient()
	if daprErr != nil {
		panic(daprErr)
	}

	result, err := daprClient.GetState(ctx, STATE_STORE_NAME, "values", nil)
	if err != nil {
		panic(err)
	}
	myValues := MyValues{}
	json.Unmarshal(result.Value, &myValues)

	respondWithJSON(w, http.StatusOK, myValues)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func main() {
	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080"
	}

	r := mux.NewRouter()

	r.HandleFunc("/", readValues).Methods("GET")

	// Add handlers for readiness and liveness endpoints
	r.HandleFunc("/health/{endpoint:readiness|liveness}", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	// Start the server; this is a blocking call
	err := http.ListenAndServe(":"+appPort, r)
	if err != http.ErrServerClosed {
		log.Panic(err)
	}
}
