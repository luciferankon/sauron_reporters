package notifier

import (
	"errors"
	"github.com/step/sauron_reporters/pkg/notifierClient"
	"log"
	"os"
	"testing"
)

type mockSlackNotifierClient struct {
	WantErr bool
}

func (m mockSlackNotifierClient) GetUsers() ([]notifierClient.User, error) {
	if m.WantErr == true {
		return []notifierClient.User{}, errors.New("unable to get users")
	}
	return []notifierClient.User{
		{
			RealName: "Ankon",
			Name:     "Ankon",
			ID:       "someID",
		},
	}, nil
}

func (m mockSlackNotifierClient) OpenIMChannel(user string) (bool, bool, string, error) {
	if user == "notValidID" {
		return false, false, "", errors.New("not able to open direct channel")
	}
	return true, true, "channelID", nil
}

func (m mockSlackNotifierClient) SendMessage(channelID string, message string) (string, string, string, error) {
	if channelID == "validChannelID" {
		return "", "", "", nil
	}
	return "", "", "", errors.New("unable to send message")
}

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
	_, err := os.OpenFile("fileWithoutPermission", os.O_CREATE, 0111)
	if err != nil {
		log.Fatal(err)
	}
}

func TestSlackNotifier_GetUserID(t *testing.T) {
	type fields struct {
		Logger           *log.Logger
		UserNameFilePath string
		SlackClient      notifierClient.NotifierClient
	}
	type args struct {
		recipient string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr error
	}{
		{
			name: "should return userId for valid recipient",
			fields: fields{
				Logger:           nil,
				UserNameFilePath: "test_usernames.json",
				SlackClient: mockSlackNotifierClient{
					WantErr: false,
				},
			},
			args: args{
				recipient: "Ankon",
			},
			want: "someID",
		},
		{
			name: "should return empty when recipient not found",
			fields: fields{
				Logger:           nil,
				UserNameFilePath: "test_usernames.json",
				SlackClient: mockSlackNotifierClient{
					WantErr: false,
				},
			},
			args: args{
				recipient: "Tilak",
			},
			want: "",
		},
		{
			name: "should return error when there is some problem fetching users",
			fields: fields{
				Logger:           nil,
				UserNameFilePath: "test_usernames.json",
				SlackClient: mockSlackNotifierClient{
					WantErr: true,
				},
			},
			args: args{
				recipient: "",
			},
			want:    "",
			wantErr: errors.New("unable to get users"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sn := SlackNotifier{
				Logger:           tt.fields.Logger,
				UserNameFilePath: tt.fields.UserNameFilePath,
				SlackClient:      tt.fields.SlackClient,
			}
			got, err := sn.GetUserID(tt.args.recipient)
			if err != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("GetMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetUserID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlackNotifier_GetChannelID(t *testing.T) {
	type fields struct {
		Logger           *log.Logger
		UserNameFilePath string
		SlackClient      notifierClient.NotifierClient
	}
	type args struct {
		userID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr error
	}{
		{
			name: "should return channel when user id is valid",
			fields: fields{
				Logger:           nil,
				UserNameFilePath: "",
				SlackClient:      mockSlackNotifierClient{WantErr: false},
			},
			args: args{
				userID: "testId",
			},
			want:    "channelID",
			wantErr: nil,
		},
		{
			name: "should return error when userID is not valid",
			fields: fields{
				Logger:           nil,
				UserNameFilePath: "",
				SlackClient:      mockSlackNotifierClient{WantErr: false},
			},
			args: args{
				userID: "notValidID",
			},
			want:    "",
			wantErr: errors.New("not able to open direct channel"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sn := SlackNotifier{
				Logger:           tt.fields.Logger,
				UserNameFilePath: tt.fields.UserNameFilePath,
				SlackClient:      tt.fields.SlackClient,
			}
			got, err := sn.GetChannelID(tt.args.userID)
			if err != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("GetMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetChannelID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlackNotifier_SendMessage(t *testing.T) {
	type fields struct {
		Logger           *log.Logger
		UserNameFilePath string
		SlackClient      notifierClient.NotifierClient
	}
	type args struct {
		channelID string
		message   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr error
	}{
		{
			name: "should return true and no error if message is sent",
			fields: fields{
				Logger:           nil,
				UserNameFilePath: "",
				SlackClient:      mockSlackNotifierClient{WantErr: false},
			},
			args: args{
				channelID: "validChannelID",
				message:   "someMessage",
			},
			want:    true,
			wantErr: nil,
		},
		{
			name: "should return false and error if message is not sent",
			fields: fields{
				Logger:           nil,
				UserNameFilePath: "",
				SlackClient:      mockSlackNotifierClient{WantErr: false},
			},
			args: args{
				channelID: "invalidChannelID",
				message:   "someMessage",
			},
			want:    false,
			wantErr: errors.New("unable to send message"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sn := SlackNotifier{
				Logger:           tt.fields.Logger,
				UserNameFilePath: tt.fields.UserNameFilePath,
				SlackClient:      tt.fields.SlackClient,
			}
			got, err := sn.SendMessage(tt.args.channelID, tt.args.message)
			if (err != nil) && err.Error() != tt.wantErr.Error() {
				t.Errorf("SendMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SendMessage() got = %v, want %v", got, tt.want)
			}
		})
	}
}
