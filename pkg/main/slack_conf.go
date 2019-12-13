package main

import (
	"flag"
	"github.com/step/sauron_reporters/pkg/notifierClient"
	"log"

	"github.com/nlopes/slack"
	"github.com/step/sauron_reporters/pkg/notifier"
)

var userNameFilePath string

func init() {
	flag.StringVar(&userNameFilePath, "user-name-file-path", "pkg/notifier/usernames.json", "`path` where user names data is stored")
}

func getSlackClient(logger *log.Logger, apiKey string) notifier.Notifier {
	return notifier.SlackNotifier{
		Logger:           logger,
		UserNameFilePath: userNameFilePath,
		SlackClient:      notifierClient.SlackNotifierClient{Client:slack.New(apiKey)},
	}
}
