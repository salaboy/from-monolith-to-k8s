package main

import (
	"testing"
	"time"
)

func TestConferenceBuilder(t *testing.T) {
	store := ConferenceStore{}
	got := store.read()
	want := []Conference{
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

	if got[0] != want[0] || got[1] != want[1]  {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

//func TestConferencesHandler(t *testing.T) {
//	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
//	// pass 'nil' as the third parameter.
//	req, err := http.NewRequest("GET", "/conferences", nil)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
//	rr := httptest.NewRecorder()
//	handler := http.HandlerFunc(ConferencesHandler)
//
//	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
//	// directly and pass in our Request and ResponseRecorder.
//	handler.ServeHTTP(rr, req)
//
//	// Check the status code is what we expect.
//	if status := rr.Code; status != http.StatusOK {
//		t.Errorf("handler returned wrong status code: got %v want %v",
//			status, http.StatusOK)
//	}
//
//	// Check the status code is what we expect.
//	if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
//		t.Errorf("handler returned wrong Content-Type : got %v want %v",
//			contentType, "application/json")
//	}
//
//	// Check the response body is what we expect.
//	store := ConferenceStore{}
//	conferences := store.read()
//	jsonConference, _ := json.Marshal(&conferences)
//	expected := jsonConference
//	if rr.Body.String() != string(expected) {
//		t.Errorf("handler returned unexpected body: got %v want %v",
//			rr.Body.String(), expected)
//	}
//}
//
//func TestHealthCheckHandler(t *testing.T) {
//	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
//	// pass 'nil' as the third parameter.
//	req, err := http.NewRequest("GET", "/health", nil)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
//	rr := httptest.NewRecorder()
//	handler := http.HandlerFunc(HealthCheckHandler)
//
//	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
//	// directly and pass in our Request and ResponseRecorder.
//	handler.ServeHTTP(rr, req)
//
//	// Check the status code is what we expect.
//	if status := rr.Code; status != http.StatusOK {
//		t.Errorf("handler returned wrong status code: got %v want %v",
//			status, http.StatusOK)
//	}
//
//	// Check the response body is what we expect.
//	expected := `{"alive": true}`
//	if rr.Body.String() != expected {
//		t.Errorf("handler returned unexpected body: got %v want %v",
//			rr.Body.String(), expected)
//	}
//}
