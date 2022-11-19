package main

import (
	"context"
	"encoding/json"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/package-url/packageurl-go"
	"os"

	cdevents "github.com/cdevents/sdk-go/pkg/api"
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

type InvolvedObject struct {
	ApiVersion string
	Kind       string
	Name       string
	Namespace  string
	Uid        string
}

type APIServerSourceEventData struct {
	Kind           string
	Message        string
	InvolvedObject InvolvedObject
}




var cdEventsEndpoint = os.Getenv("CDEVENTS_ENDPOINT_URL")


var logger = log.NewLogfmtLogger(os.Stdout)

func CheckError(err error) {
	if err != nil {
		level.Error(logger).Log("error", err)
		panic(err)
	}
}
var ceClient cloudevents.Client

func main() {

	// The default client is HTTP.
	ceClient, err := cloudevents.NewClientHTTP()
	if err != nil {
		level.Error(logger).Log("failed to create client, %v", err)
		return
	}
	level.Error(logger).Log(ceClient.StartReceiver(context.Background(), receive))
}


func receive(event cloudevents.Event) {

	// Deal with Cloud Events coming from APIServer Source
	if event.Type() == "dev.knative.apiserver.resource.add" {

		// Deal with Deployments

		// First parse the APIServerSourceEventData to find the involved object
		data := APIServerSourceEventData{}

		err := json.Unmarshal(event.Data(), &data)
		if err != nil {
			level.Error(logger).Log("Failed to parse APIServerSourceEventData ", err)
		}

		// Only transform deployment events
		if data.InvolvedObject.Kind == "Deployment" {
			level.Debug(logger).Log("Deployment event found ... from: ",event.Type() )
			cdevent, _ := cdevents.NewServiceDeployedEvent()
			cdevent.SetTimestamp(event.Time())
			cdevent.SetSubjectSource("ApiServerSource")
			cdevent.SetSubjectId(data.InvolvedObject.Name)
			cdevent.SetSource(event.Source())
			environment := cdevents.Reference{}
			environment.Id = data.InvolvedObject.Namespace
			environment.Source = event.Source()
			cdevent.SetSubjectEnvironment(environment)
			cdevent.SetCustomData("application/json", event)

			cdevent.SetSubjectArtifactId(packageurl.NewPackageURL("pkg", data.InvolvedObject.Namespace, data.InvolvedObject.Name, data.InvolvedObject.Uid, nil, "").ToString())

			ce, err := cdevents.AsCloudEvent(cdevent)
			if err != nil {
				level.Debug(logger).Log("Failed to transform a CDEvent into a CloudEvent ", err)
				return
			}

			level.Debug(logger).Log("Sending CloudEvent to: ",cdEventsEndpoint )
			ctx := cloudevents.ContextWithTarget(context.Background(), cdEventsEndpoint)
			ctx = cloudevents.WithEncodingBinary(ctx)
			ceSenderClient, err := cloudevents.NewClientHTTP()
			if err != nil {
				level.Error(logger).Log("failed to create client, %v", err)
				return
			}
			if result := ceSenderClient.Send(ctx, *ce); cloudevents.IsUndelivered(result) {
				level.Debug(logger).Log("failed to send:", result)
			}

		}

	}


}


