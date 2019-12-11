package writer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	st "github.com/step/saurontypes"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDbWriter struct {
	URI    string
	DB     string
	Table  string
	Logger *log.Logger
}

func generateDBReport(report string, event map[string]interface{}) st.DBReport {
	var reportJSON st.Report
	var results st.Results
	var testResult st.TestResult
	json.Unmarshal([]byte(report), &reportJSON)
	json.Unmarshal([]byte(reportJSON.Results), &results)
	json.Unmarshal([]byte(results.Results), &testResult)

	return st.DBReport{
		Job:     reportJSON.Job,
		Result:  testResult,
		FlowID:  fmt.Sprintf("%v", event["flowID"]),
		Project: fmt.Sprintf("%v", event["project"]),
		Pusher:  fmt.Sprintf("%v", event["pusherID"]),
		Time:    fmt.Sprintf("%v", event["timestamp"]),
	}
}

func (mdbwriter MongoDbWriter) Write(events map[string]interface{}) {
	clientOptions := options.Client().ApplyURI(mdbwriter.URI)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	eventsCollection := client.Database(mdbwriter.DB).Collection(mdbwriter.Table)
	var docs []interface{}
	details := fmt.Sprintf("%v", events["details"])
	report := details[strings.IndexByte(details, '{'):]
	dbReport := generateDBReport(report, events)
	docs = append(docs, dbReport)

	eventsCollection.InsertMany(context.TODO(), docs)
	mdbwriter.logWrite(dbReport.FlowID)
}
