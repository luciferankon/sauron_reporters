package main

import (
	"flag"
	"github.com/go-redis/redis"
	"time"
	"github.com/step/sauron_reporters/pkg/event_listner"
)

func main() {
	flag.Parse()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	r := make(chan bool, 100)
	stop := make(chan bool)

	listner := eventlistner.NewListner(*redisClient)

	go listner.Start(redisClient, "mystream", r, stop)

	for range r {
		time.Sleep(time.Millisecond * 100)
	}
}