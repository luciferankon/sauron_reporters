package main

import (
	"flag"

	"github.com/step/sauron_reporters/pkg/notifier"
)

var apiKey string

func init() {
	flag.StringVar(&apiKey, "api-key", "xoxb-681535316417-677029338466-FuflSk6ovXQPedTeXvhXpBqm", "api key of the slack app")
}

func getSlackClient() notifier.Notifier {
	return notifier.SlackNotifier{
		ApiKey: apiKey,
	}
}
