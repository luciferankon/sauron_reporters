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

func GenerateDBReport(report string, event map[string]interface{}) (st.DBReport, error) {
	var reportJSON st.Report
	var results st.Results
	var testResult st.TestResult
	err := json.Unmarshal([]byte(report), &reportJSON)
	if err != nil {
		return st.DBReport{},err
	}

	err = json.Unmarshal([]byte(reportJSON.Results), &results)
	if err != nil {
		return st.DBReport{},err
	}

	err = json.Unmarshal([]byte(results.Results), &testResult)
	if err != nil {
		return st.DBReport{},err
	}

	return st.DBReport{
		Job:     reportJSON.Job,
		Result:  testResult,
		FlowID:  fmt.Sprintf("%v", event["flowID"]),
		Project: fmt.Sprintf("%v", event["project"]),
		Pusher:  fmt.Sprintf("%v", event["pusherID"]),
		Time:    fmt.Sprintf("%v", event["timestamp"]),
	} , nil
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
	dbReport, err := GenerateDBReport(report, events)
	if err != nil {
		mdbwriter.logError("Unable to generate report", err)
	}
	docs = append(docs, dbReport)

	eventsCollection.InsertMany(context.TODO(), docs)
	mdbwriter.logWrite(dbReport.FlowID)
}
