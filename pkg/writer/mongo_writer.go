package writer

import (
	"context"
	"encoding/json"
	"fmt"
	st "github.com/step/saurontypes"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"strings"
)

type MongoDbWriter struct {
	Client *mongo.Client
	DB     string
	Table  string
	Logger *log.Logger
}

func GenerateDBReport(report string, event map[string]interface{}) (interface{}, error) {
	var reportJSON st.Report
	err := json.Unmarshal([]byte(report), &reportJSON)
	if err != nil {
		return st.DBTestReport{}, err
	}

	var results st.Results
	err = json.Unmarshal([]byte(reportJSON.Results), &results)
	if err != nil {
		return st.DBTestReport{}, err
	}

	if reportJSON.Job == "test" {
		return GenerateDBReportForTest(results.Results, event)
	}
	return GenerateDBReportForLint(results.Results, event)
}

func GenerateDBReportForTest(results string, event map[string]interface{}) (interface{}, error) {
	var testResult st.TestResult

	err := json.Unmarshal([]byte(results), &testResult)
	if err != nil {
		return st.DBTestReport{}, err
	}

	return st.DBTestReport{
		Job:     "test",
		Result:  testResult,
		FlowID:  fmt.Sprintf("%v", event["flowID"]),
		Project: fmt.Sprintf("%v", event["project"]),
		Pusher:  fmt.Sprintf("%v", event["pusherID"]),
		Time:    fmt.Sprintf("%v", event["timestamp"]),
		SHA:     fmt.Sprintf("%v", event["sha"]),
	}, nil
}

func GenerateDBReportForLint(result string, event map[string]interface{}) (interface{}, error) {
	var lintResult []st.LintResult
	err := json.Unmarshal([]byte(result), &lintResult)
	if err != nil {
		return st.DBLintReport{}, err
	}

	return st.DBLintReport{
		Job:     "lint",
		Result:  lintResult,
		FlowID:  fmt.Sprintf("%v", event["flowID"]),
		Project: fmt.Sprintf("%v", event["project"]),
		Pusher:  fmt.Sprintf("%v", event["pusherID"]),
		Time:    fmt.Sprintf("%v", event["timestamp"]),
		SHA:     fmt.Sprintf("%v", event["sha"]),
	}, nil
}

func (mdbwriter MongoDbWriter) Write(events map[string]interface{}) {

	// clientOptions := options.Client().ApplyURI(mdbwriter.URI)
	// client, err := mongo.Connect(context.TODO(), clientOptions)

	err := mdbwriter.Client.Ping(context.TODO(), nil)

	if err != nil {
		mdbwriter.logError("Unable to ping due to => ", err)
		return
	}

	eventsCollection := mdbwriter.Client.Database(mdbwriter.DB).Collection(mdbwriter.Table)
	var docs []interface{}
	details := fmt.Sprintf("%v", events["details"])
	report := details[strings.IndexByte(details, '{'):]
	dbReport, err := GenerateDBReport(report, events)
	if err != nil {
		mdbwriter.logError("Unable to generate report", err)
		return
	}
	docs = append(docs, dbReport)

	_, err = eventsCollection.InsertMany(context.TODO(), docs)
	if err != nil {
		mdbwriter.logError("Unable to insert data due to => ", err)
		return
	}
	mdbwriter.logWrite(fmt.Sprintf("%v", events["flowID"]))
}
