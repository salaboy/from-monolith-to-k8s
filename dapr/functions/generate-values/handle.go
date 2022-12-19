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

	"github.com/google/uuid"
)

var stateStoreName = "statestore"
var daprPort = os.Getenv("DAPR_HTTP_PORT")
var stateStoreUrl = "http://localhost:" + daprPort + "/v1.0/state/" + stateStoreName

type MyObject struct {
	Key   string
	Value string
}

// Handle an HTTP Request.
func Handle(ctx context.Context, res http.ResponseWriter, req *http.Request) {

	fmt.Println("Received request")
	id := uuid.New()
	myObj := MyObject{
		Key:   id.String(),
		Value: strconv.Itoa(rand.Intn(100)),
	}

	myObjs := []MyObject{myObj}

	json_data, err := json.Marshal(myObjs)

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
