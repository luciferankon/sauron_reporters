package writer

import (
	"errors"
	"github.com/step/saurontypes"
	"reflect"
	"testing"
)

func TestGenerateDBReport(t *testing.T) {
	type args struct {
		report string
		event  map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    saurontypes.DBReport
		wantErr error
	}{
		{
			name: "should return db-report for valid report and event",
			args: args{
				report: `{"job":"test","result":"{\"result.json\":\"{\\\"total\\\":10,\\\"passed\\\":[],\\\"failed\\\":[],\\\"pending\\\":[]}\"}"}`,
				event: map[string]interface{}{
					"flowID":    "123",
					"project":   "head",
					"pusherID":  "luciferankon",
					"timestamp": "2019-12-12T07:03:28.228Z",
				},
			},
			want: saurontypes.DBReport{
				Job: "test",
				Result: saurontypes.TestResult{
					Total:   10,
					Passed:  []saurontypes.TestReport{},
					Failed:  []saurontypes.TestReport{},
					Pending: []saurontypes.TestReport{},
				},
				FlowID:  "123",
				Project: "head",
				Pusher:  "luciferankon",
				Time:    "2019-12-12T07:03:28.228Z",
			},
			wantErr: nil,
		},
		{
			name: "should return unmarshal error for invalid job format",
			args: args{
				report: `{"job":"test,"result":"{\"result.json\":\"{\\\"total\\\":10,\\\"passed\\\":[],\\\"failed\\\":[],\\\"pending\\\":[]}\"}"}`,
				event: map[string]interface{}{
					"flowID":    "123",
					"project":   "head",
					"pusherID":  "luciferankon",
					"timestamp": "2019-12-12T07:03:28.228Z",
				},
			},
			want:    saurontypes.DBReport{},
			wantErr: errors.New("invalid character 'r' after object key:value pair"),
		},
		{
			name: "should return unmarshal error for invalid result.json",
			args: args{
				report: `{"job":"test","result":"{result.json\":\"{\\\"total\\\":10,\\\"passed\\\":[],\\\"failed\\\":[],\\\"pending\\\":[]}\"}"}`,
				event: map[string]interface{}{
					"flowID":    "123",
					"project":   "head",
					"pusherID":  "luciferankon",
					"timestamp": "2019-12-12T07:03:28.228Z",
				},
			},
			want:    saurontypes.DBReport{},
			wantErr: errors.New("invalid character 'r' looking for beginning of object key string"),
		},
		{
			name: "should return unmarshal error for invalid total",
			args: args{
				report: `{"job":"test","result":"{\"result.json\":\"{total\\\":10,\\\"passed\\\":[],\\\"failed\\\":[],\\\"pending\\\":[]}\"}"}`,
				event: map[string]interface{}{
					"flowID":    "123",
					"project":   "head",
					"pusherID":  "luciferankon",
					"timestamp": "2019-12-12T07:03:28.228Z",
				},
			},
			want:    saurontypes.DBReport{},
			wantErr: errors.New("invalid character 't' looking for beginning of object key string"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateDBReport(tt.args.report, tt.args.event)
			if (err != nil) && err.Error() != tt.wantErr.Error() {
				t.Errorf("GetMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateDBReport() = %v, want %v", got, tt.want)
			}
		})
	}
}
