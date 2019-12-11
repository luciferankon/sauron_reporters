package main

import (
	"log"

	"github.com/step/sauron_reporters/pkg/notifier"
)

func getSlackClient(logger *log.Logger, apiKey string) notifier.Notifier {
	return notifier.SlackNotifier{
		ApiKey: apiKey,
		Logger: logger,
	}
}
