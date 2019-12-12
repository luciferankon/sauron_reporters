package notifier

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/nlopes/slack"
	st "github.com/step/saurontypes"
)

type SlackNotifier struct {
	ApiKey           string
	Logger           *log.Logger
	UserNameFilePath string
}

type Users struct {
	Users []User `json:"users"`
}

type User struct {
	GithubUserName string `json:"githubUserName"`
	SlackUserName  string `json:"slackUserName"`
}

func GetMessage(report string) (string, error) {
	var reportJSON st.Report
	var results st.Results
	var testResult st.TestResult
	err := json.Unmarshal([]byte(report), &reportJSON)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal([]byte(reportJSON.Results), &results)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal([]byte(results.Results), &testResult)
	if err != nil {
		return "", err
	}

	message := fmt.Sprintf("```\nTest Results\nTotal => %d\nPassed => %d\nFailed => %d\n```", testResult.Total, len(testResult.Passed), len(testResult.Failed))
	return message, nil
}

func GetUserName(githubUserName interface{}, userDataFilePath string) (string, error) {
	jsonFile, err := os.Open(userDataFilePath)
	if err != nil {
		return "", err
	}

	jsonInBytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return "", err
	}

	var users Users
	err = json.Unmarshal(jsonInBytes, &users)
	if err != nil {
		return "", err
	}

	for _, user := range users.Users {
		if user.GithubUserName == githubUserName {
			return user.SlackUserName, nil
		}
	}
	return "", errors.New("username not found")
}

func (sn SlackNotifier) GetAPI() *slack.Client {
	return slack.New(sn.ApiKey)
}

func (sn SlackNotifier) GetUserID(recipient string, api *slack.Client) string {
	users, err := api.GetUsers()
	if err != nil {
		log.Fatalf("Not able to get users due to ==> %s", err)
	}

	for _, user := range users {
		if (user.RealName == recipient) || (user.Name == recipient) {
			return user.ID
		}
	}
	return ""
}

func (sn SlackNotifier) GetChannelID(api *slack.Client, userID string) string {
	_, _, channelID, err := api.OpenIMChannel(userID)
	if err != nil {
		log.Fatalf("Not able to open direct channel due to ==> %s", err)
	}
	return channelID
}

func (sn SlackNotifier) SendMessage(channelID, message string, api *slack.Client) {
	m := slack.MsgOptionText(message, true)
	_, _, _, err := api.SendMessage(channelID, m)
	if err != nil {
		log.Fatalf("Not able to send message due to ==> %s", err)
	}
}

func (sn SlackNotifier) Notify(events map[string]interface{}) {
	recipient, err := GetUserName(events["pusherID"], sn.UserNameFilePath)
	if err != nil {
		sn.logError("Unable to find username", err)
		return
	}

	details := fmt.Sprintf("%v", events["details"])
	report := details[strings.IndexByte(details, '{'):]

	message, err := GetMessage(report)
	if err != nil {
		sn.logError("Unable to generate message", err)
		return
	}

	api := sn.GetAPI()
	userID := sn.GetUserID(recipient, api)
	channelID := sn.GetChannelID(api, userID)

	sn.SendMessage(channelID, message, api)
	sn.logNotification(message, recipient)
}
