package main

import (
	"flag"
	"log"

	"github.com/step/sauron_reporters/pkg/notifier"
)

var userNameFilePath string

func init() {
	flag.StringVar(&userNameFilePath, "user-name-file-path", "pkg/notifier/usernames.json", "`path` where user names data is stored")
}

func getSlackClient(logger *log.Logger, apiKey string) notifier.Notifier {
	return notifier.SlackNotifier{
		ApiKey:           apiKey,
		Logger:           logger,
		UserNameFilePath: userNameFilePath,
	}
}
