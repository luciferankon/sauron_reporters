package dbwriter

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

func generateDBReport(report string, event redis.XMessage) DBReport {
	var reportJSON Report
	var results Results
	var testResult TestResult
	json.Unmarshal([]byte(report), &reportJSON)
	json.Unmarshal([]byte(reportJSON.Results), &results)
	json.Unmarshal([]byte(results.Results), &testResult)

	return DBReport{
		Job:    reportJSON.Job,
		Result: testResult,
		FlowID: fmt.Sprintf("%v",event.Values["flowID"]),
		Project: fmt.Sprintf("%v",event.Values["project"]),
		Pusher: fmt.Sprintf("%v",event.Values["pusherID"]),
		Time: fmt.Sprintf("%v",event.Values["timestamp"]),
	}
}

func Write(events []redis.XMessage) {
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
	for _, event := range events {
		details := fmt.Sprintf("%v", event.Values["details"])
		report := details[strings.IndexByte(details, '{'):]
		dbReport := generateDBReport(report, event)
		docs = append(docs, dbReport)
	}

	eventsCollection.InsertMany(context.TODO(), docs)
}
