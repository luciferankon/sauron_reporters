package notifier

import (
	"errors"
	"log"
	"os"
	"testing"
)

func TestGetMessage(t *testing.T) {
	type args struct {
		report string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "should return message for valid report",
			args: args{
				report: `{"job":"test","result":"{\"result.json\":\"{\\\"total\\\":10,\\\"passed\\\":[],\\\"failed\\\":[],\\\"pending\\\":[]}\"}"}`,
			},
			want:    "```\nTest Results\nTotal => 10\nPassed => 0\nFailed => 0\n```",
			wantErr: false,
		},
		{
			name: "should return error for invalid job",
			args: args{
				report: `{"job":"test,"result":"{\"result.json\":\"{\\\"total\\\":10,\\\"passed\\\":[],\\\"failed\\\":[],\\\"pending\\\":[]}\"}"}`,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "should return error for invalid result.json",
			args: args{
				report: `{"job":"test","result":"{result.json\":\"{\\\"total\\\":10,\\\"passed\\\":[],\\\"failed\\\":[],\\\"pending\\\":[]}\"}"}`,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "should return error for invalid total count",
			args: args{
				report: `{"job":"test","result":"{\"result.json\":\"{total\\\":10,\\\"passed\\\":[],\\\"failed\\\":[],\\\"pending\\\":[]}\"}"}`,
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetMessage(tt.args.report)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetMessage() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetUserName(t *testing.T) {
	setup()
	type args struct {
		githubUserName   interface{}
		userDataFilePath string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr error
	}{
		{
			name: "Should return slack username associated to given github user name",
			args: args{
				githubUserName:   "luciferankon",
				userDataFilePath: "test_usernames.json",
			},
			want:    "Ankon",
			wantErr: nil,
		},
		{
			name: "Should return error if the github user is not associated",
			args: args{
				githubUserName:   "luciferanko",
				userDataFilePath: "test_usernames.json",
			},
			want:    "",
			wantErr: errors.New("username not found"),
		},
		{
			name: "Should return error if the file is not readable",
			args: args{
				githubUserName:   "luciferankon",
				userDataFilePath: "fileWithoutPermission",
			},
			want:    "",
			wantErr: errors.New("open fileWithoutPermission: permission denied"),
		},
		{
			name: "Should return error if the file is not valid json",
			args: args{
				githubUserName:   "luciferankon",
				userDataFilePath: "invalid_json.json",
			},
			want:    "",
			wantErr: errors.New("invalid character 'x' looking for beginning of object key string"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUserName(tt.args.githubUserName, tt.args.userDataFilePath)
			if (err != nil) && err.Error() != tt.wantErr.Error() {
				t.Errorf("GetUserName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetUserName() got = %v, want %v", got, tt.want)
			}
		})
	}
	teardown()
}

func teardown() {
	err := os.Remove("fileWithoutPermission")
	if err != nil {
		log.Fatal(err)
	}
}

func setup() {
	_, err := os.OpenFile("fileWithoutPermission",os.O_CREATE, 0111)
	if err != nil {
		log.Fatal(err)
	}
}