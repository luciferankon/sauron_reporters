package main

import (
	"flag"
	"time"

	"github.com/step/sauron_reporters/pkg/eventlistener"
)

func main() {
	flag.Parse()

	r := make(chan bool, 100)
	stop := make(chan bool)
	redisClient := getRedisClient()
	mongoWriter := getMongoWriter()
	slackClient := getSlackClient()

	listener := eventlistener.NewListner(redisClient, mongoWriter, slackClient)

	go listener.Start("eventHub", r, stop)

	for range r {
		time.Sleep(time.Millisecond * 100)
	}
}
