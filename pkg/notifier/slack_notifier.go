package notifier

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/nlopes/slack"
	st "github.com/step/saurontypes"
)

type SlackNotifier struct {
	ApiKey string
}

type Users struct {
	Users []User `json:"users"`
}

type User struct {
	GithubUserName string `json:"githubUserName"`
	SlackUserName  string `json:"slackUserName"`
}

func getMessage(report string, events map[string]interface{}) string {
	var reportJSON st.Report
	var results st.Results
	var testResult st.TestResult
	json.Unmarshal([]byte(report), &reportJSON)
	json.Unmarshal([]byte(reportJSON.Results), &results)
	json.Unmarshal([]byte(results.Results), &testResult)

	message := fmt.Sprintf("```Test Results\nTotal => %d\nPassed => %d\nFailed => %d\n```", testResult.Total, len(testResult.Passed), len(testResult.Failed))
	return message
}

func getUserName(githubUserName interface{}) string {
	jsonFile, _ := os.Open("pkg/notifier/usernames.json")
	jsonInBytes, _ := ioutil.ReadAll(jsonFile)

	var users Users
	json.Unmarshal(jsonInBytes, &users)
	for _, user := range users.Users {
		if user.GithubUserName == githubUserName {
			return user.SlackUserName
		}
	}
	return ""
}

func (sn SlackNotifier) getAPI() *slack.Client {
	return slack.New(sn.ApiKey)
}

func (sn SlackNotifier) getUserID(recipient string, api *slack.Client) string {
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

func (sn SlackNotifier) getChannelID(api *slack.Client, userID string) string {
	_, _, channelID, err := api.OpenIMChannel(userID)
	if err != nil {
		log.Fatalf("Not able to open direct channel due to ==> %s", err)
	}
	return channelID
}

func (sn SlackNotifier) sendMessage(channelID, message string, api *slack.Client) {
	m := slack.MsgOptionText(message, true)
	_, _, _, err := api.SendMessage(channelID, m)
	if err != nil {
		log.Fatalf("Not able to send message due to ==> %s", err)
	}
}

func (sn SlackNotifier) Notify(events map[string]interface{}) {
	recipient := getUserName(events["pusherID"])
	details := fmt.Sprintf("%v", events["details"])
	report := details[strings.IndexByte(details, '{'):]
	message := getMessage(report, events)

	api := sn.getAPI()
	userID := sn.getUserID(recipient, api)
	channelID := sn.getChannelID(api, userID)

	sn.sendMessage(channelID, message, api)
}
