package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/gorilla/mux"
)

type Input struct {
	ID     string
	Value  string
	Stored bool
}

type Result struct {
	ID        string
	Input     string
	Output    string
	Processed bool
}

type Results struct {
	Results []Result
}

var results Results

var inputs Inputs

type Inputs struct {
	Inputs []Input
}

type MyValues struct {
	Values []string
}

type MyObject struct {
	Key   string
	Value MyValues
}

var STATE_STORE_NAME = "statestore"
var daprClient dapr.Client
var daprErr error

func main() {

	daprClient, daprErr = dapr.NewClient()
	if daprErr != nil {
		panic(daprErr)
	}

	r := mux.NewRouter()

	r.HandleFunc("/info", InfoHandler).Methods("GET")
	r.HandleFunc("/avg", AverageHandler).Methods("GET")
	r.HandleFunc("/store", StoreHandler).Methods("POST")
	r.HandleFunc("/values", GetValuesHandler).Methods("GET")

	r.HandleFunc("/clear", ClearHandler).Methods("DELETE")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(os.Getenv("KO_DATA_PATH"))))
	log.Printf("Dapr+Knative Functions app Started in port 8080!")

	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8080", nil))
	defer daprClient.Close()

}

func InfoHandler(writer http.ResponseWriter, request *http.Request) {
	respondWithJSON(writer, http.StatusOK, "{ 'app': 'OK' }")
}

func AverageHandler(writer http.ResponseWriter, request *http.Request) {
	resp, err := http.Get("http://avg.default.svc.cluster.local")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	writer.Header().Set("Content-Type", "text/html")
	writer.WriteHeader(http.StatusOK)
	writer.Write(body)
}

// func StoreHandler shoudl:
//
//	Check in the state store if there is a key called `values`
//	- If it exist it should retrieve the content and add the new value to its Values
//	- If it doesn't exist it should create the key with an array of values which includes a single element (the value provided in this request)
func StoreHandler(writer http.ResponseWriter, request *http.Request) {

	value := request.URL.Query().Get("value")

	ctx := context.Background()

	result, _ := daprClient.GetState(ctx, STATE_STORE_NAME, "values", nil)
	myValues := MyValues{}
	if result.Value != nil {
		json.Unmarshal(result.Value, &myValues)
	}

	if myValues.Values == nil || len(myValues.Values) == 0 {
		myValues.Values = []string{value}
	} else {
		myValues.Values = append(myValues.Values, value)
	}

	jsonData, err := json.Marshal(myValues)

	err = daprClient.SaveState(ctx, STATE_STORE_NAME, "values", jsonData, nil)
	if err != nil {
		panic(err)
	}

	respondWithJSON(writer, http.StatusOK, jsonData)

}

func ClearHandler(writer http.ResponseWriter, request *http.Request) {

	// deleted, err := client.Del("values").Result()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// respondWithJSON(writer, http.StatusOK, deleted)
}

func GetValuesHandler(writer http.ResponseWriter, request *http.Request) {

	ctx := context.Background()

	result, err := daprClient.GetState(ctx, STATE_STORE_NAME, "values", nil)
	if err != nil {
		panic(err)
		respondWithJSON(writer, http.StatusOK, "[]")
		return
	}
	myValues := MyValues{}
	json.Unmarshal(result.Value, &myValues)

	respondWithJSON(writer, http.StatusOK, myValues.Values)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
