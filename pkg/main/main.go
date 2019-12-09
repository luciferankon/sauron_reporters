package main

import (
	"flag"
	"time"

	"github.com/step/sauron_reporters/pkg/eventlistener"
	w "github.com/step/sauron_reporters/pkg/writer"
)

func main() {
	flag.Parse()

	r := make(chan bool, 100)
	stop := make(chan bool)
	redisClient := getRedisClient()

	mongoWriter := w.NewMongoWriter()
	listener := eventlistener.NewListner(redisClient, mongoWriter)

	go listener.Start("eventHub", r, stop)

	for range r {
		time.Sleep(time.Millisecond * 100)
	}
}
