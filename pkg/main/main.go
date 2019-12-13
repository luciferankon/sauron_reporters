package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/step/sauron_reporters/pkg/eventlistener"
)

func main() {
	flag.Parse()

	listenerLogger := log.New(os.Stdout, "listener ", log.LstdFlags)
	notifierLogger := log.New(os.Stdout, "notifier ", log.LstdFlags)
	writerLogger := log.New(os.Stdout, "writer ", log.LstdFlags)

	slackAuth := os.Getenv("API_KEY")

	r := make(chan bool, 100)
	stop := make(chan bool)
	redisClient := getRedisClient()
	mongoWriter := getMongoWriter(writerLogger)
	slackClient := getSlackClient(notifierLogger, slackAuth)

	listener := eventlistener.NewListner(redisClient, mongoWriter, slackClient, listenerLogger)

	go listener.Start("eventHub", r, stop)

	for range r {
		time.Sleep(time.Millisecond * 100)
	}
}
