package writer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDbWriter struct{}

type Report struct {
	Job     string `json:"job"`
	Results string `json:"result"`
}

type Results struct {
	Results string `json:"result.json"`
}

type TestResult struct {
	Total   int          `json:"total"`
	Passed  []TestReport `json:"passed"`
	Failed  []TestReport `json:"failed"`
	Pending []TestReport `json:"pending"`
}

type TestReport struct {
	Suite string `json:"suite"`
	Title string `json:"title"`
}

type DBReport struct {
	Job     string
	Result  TestResult
	FlowID  string
	Project string
	Pusher  string
	Time    string
}

func generateDBReport(report string, event map[string]interface{}) DBReport {
	var reportJSON Report
	var results Results
	var testResult TestResult
	json.Unmarshal([]byte(report), &reportJSON)
	json.Unmarshal([]byte(reportJSON.Results), &results)
	json.Unmarshal([]byte(results.Results), &testResult)

	return DBReport{
		Job:     reportJSON.Job,
		Result:  testResult,
		FlowID:  fmt.Sprintf("%v", event["flowID"]),
		Project: fmt.Sprintf("%v", event["project"]),
		Pusher:  fmt.Sprintf("%v", event["pusherID"]),
		Time:    fmt.Sprintf("%v", event["timestamp"]),
	}
}

func (mdbwriter MongoDbWriter) Write(events map[string]interface{}) {
	clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	eventsCollection := client.Database("sauron_reporters").Collection("events")
	var docs []interface{}
	details := fmt.Sprintf("%v", events["details"])
	report := details[strings.IndexByte(details, '{'):]
	dbReport := generateDBReport(report, events)
	docs = append(docs, dbReport)

	eventsCollection.InsertMany(context.TODO(), docs)
}

func NewMongoWriter() MongoDbWriter {
	return MongoDbWriter{}
}
