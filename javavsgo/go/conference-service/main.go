package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"time"
)

type Conference struct {
	Id    string    `json:"id"`
	Name  string    `json:"name"`
	Where string    `json:"where"`
	When  time.Time `json:"when"`
}

////time.Date(2022, time.July, 18, 0, 0, 0, 0, time.UTC)
//func main(){
//  // Creamos una instancia de nuestra struct conferencia
//
//	//Printout
//
//}














type ConferenceStore struct {
}

func (*ConferenceStore) read() []Conference {
	return []Conference{
		{
			Id:    "123",
			Name:  "JBCNConf",
			Where: "Barcelona, Spain",
			When:  time.Date(2022, time.July, 18, 0, 0, 0, 0, time.UTC),
		},
		{
			Id:    "456",
			Name:  "KubeCon",
			Where: "Detroit, USA",
			When:  time.Date(2022, time.October, 24, 0, 0, 0, 0, time.UTC),
		},
	}
}

//func main(){
//	store := ConferenceStore{}
//	conferences := store.read()
//	log.Println("Mis conferencia favorita son: ", conferences)
//}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/conferences", ConferencesHandler).Methods(http.MethodGet)
	//r.HandleFunc("/health", HealthCheckHandler).Methods(http.MethodGet)
	http.Handle("/", r)
	log.Println("HTTP Server started on Port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func ConferencesHandler(w http.ResponseWriter, r *http.Request) {
	store := ConferenceStore{}
	conferences := store.read()
	response, _ := json.Marshal(&conferences)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}


func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, `{"alive": true}`)
}




