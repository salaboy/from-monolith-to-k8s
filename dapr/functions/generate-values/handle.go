package function

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
)

var stateStoreName = "statestore"
var daprPort = os.Getenv("DAPR_HTTP_PORT")
var stateStoreUrl = "http://localhost:" + daprPort + "/v1.0/state/" + stateStoreName

// Handle an HTTP Request.
func Handle(ctx context.Context, res http.ResponseWriter, req *http.Request) {

	fmt.Println("Received request")
	valueString := strconv.Itoa(rand.Intn(100))
	value := map[string]string{"key": "random-" + valueString, "value": valueString}
	objects := []interface{}{value}
	json_data, err := json.Marshal(objects)

	fmt.Println("Json data: ", string(json_data))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Storing state to: ", stateStoreUrl)

	resp, err := http.Post(stateStoreUrl, "application/json",
		bytes.NewBuffer(json_data))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	fmt.Println("Storing state code: ", resp.StatusCode, http.StatusText(resp.StatusCode))

	fmt.Fprintln(res, "OK")

}
