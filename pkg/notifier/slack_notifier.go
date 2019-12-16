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
	err := json.Unmarshal([]byte(report), &reportJSON)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal([]byte(reportJSON.Results), &results)
	if err != nil {
		return "", err
	}

	if reportJSON.Job == "test" {
		return generateMessageForTest(results.Results)
	}
	return generateMessageForLint(results.Results)
}

func generateMessageForTest(results string) (string, error) {
	var testResult st.TestResult
	err := json.Unmarshal([]byte(results), &testResult)
	if err != nil {
		return "", err
	}
	message := fmt.Sprintf("```\nTest Results\nTotal => %d\nPassed => %d\nFailed => %d\n```", testResult.Total, len(testResult.Passed), len(testResult.Failed))
	return message, nil
}

func generateMessageForLint(results string) (string, error) {
	var lintResults []st.LintResult
	err := json.Unmarshal([]byte(results), &lintResults)
	if err != nil {
		return "", err
	}
	var message strings.Builder
	errorsCount := 0
	for _, lintResult := range lintResults {
		errorsCount += lintResult.ErrorCount
		message.WriteString("---\n")
		message.WriteString("Filename => " + lintResult.FileName + "\n")
		message.WriteString("Error count => " + fmt.Sprintf("%d", lintResult.ErrorCount) + "\n")
		message.WriteString("\nErrors\n======\n")

		for _, lintMessage := range lintResult.Messages {
			message.WriteString("`Error " + lintMessage.Message + " found at line " + fmt.Sprintf("%d", lintMessage.Line) + " column " + fmt.Sprintf("%d", lintMessage.Col) + "`\n")
		}

		message.WriteString("---\n\n")
	}
	message.WriteString("Total errors => " + fmt.Sprintf("%d", errorsCount) + "\n")
	return message.String(), nil
}

func GetUserName(githubUserName interface{}, userDataFilePath string) (string, error) {
	jsonFile, err := os.Open(userDataFilePath)

	
	if err != nil {
		return "", err
	}

	defer func() {
		fmt.Println("############File closing!!")
		if err := jsonFile.Close(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

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
	return "", nil
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

	isSent, err := sn.SendMessage(channelID, fmt.Sprintf("%s\nPusher => %v\n",message,events["pusherID"]))
	if !isSent && err != nil {
		sn.logError("Unable to send message due to => ", err)
		return
	}
	sn.logNotification(message, recipient)
}
