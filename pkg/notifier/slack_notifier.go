package notifier

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/step/sauron_reporters/pkg/notifierClient"
	"io/ioutil"
	"log"
	"os"
	"strings"

	st "github.com/step/saurontypes"
)

type SlackNotifier struct {
	Logger           *log.Logger
	UserNameFilePath string
	SlackClient      notifierClient.NotifierClient
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

func (sn SlackNotifier) GetUserID(recipient string) (string, error) {
	users, err := sn.SlackClient.GetUsers()
	if err != nil {
		return "", err
	}

	for _, user := range users {
		if (user.RealName == recipient) || (user.Name == recipient) {
			return user.ID, nil
		}
	}
	return "",nil
}

func (sn SlackNotifier) GetChannelID(userID string) (string, error) {
	_, _, channelID, err := sn.SlackClient.OpenIMChannel(userID)
	if err != nil {
		return "", err
	}
	return channelID, nil
}

func (sn SlackNotifier) SendMessage(channelID, message string) (bool, error) {
	_, _, _, err := sn.SlackClient.SendMessage(channelID, message)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (sn SlackNotifier) Notify(events map[string]interface{}) {
	recipient, err := GetUserName(events["pusherID"], sn.UserNameFilePath)
	if err != nil {
		sn.logError("Unable to find username => ", err)
		return
	}

	details := fmt.Sprintf("%v", events["details"])
	report := details[strings.IndexByte(details, '{'):]

	message, err := GetMessage(report)
	if err != nil {
		sn.logError("Unable to generate message due to => ", err)
		return
	}

	userID, err := sn.GetUserID(recipient)
	if err != nil {
		sn.logError("Unable to get userID due to => ", err)
		return
	}
	channelID, err := sn.GetChannelID(userID)
	if err != nil {
		sn.logError("Unable to get channelID due to => ", err)
		return
	}

	isSent, err := sn.SendMessage(channelID, message)
	if !isSent && err != nil {
		sn.logError("Unable to send message due to => ", err)
		return
	}
	sn.logNotification(message, recipient)
}
