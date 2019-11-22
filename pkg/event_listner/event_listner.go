package eventlistner

import (
	"github.com/go-redis/redis"
	"time"
	"fmt"
	"github.com/step/sauron_reporters/pkg/db_writer"
)

type Listner struct {
	StremClient redis.Client
}

func (l Listner) Start(redisClient *redis.Client, streamName string, r chan<- bool, stop <-chan bool) {
	shouldStop := false
	lastIDRead := "0"
	go func() {
		shouldStop = <-stop
	}()

	for {
		// Keep running till asked to stop
		if shouldStop {
			break
		}
		readArgs := &redis.XReadArgs{
			Streams: []string{streamName, lastIDRead},
		}
		// Take all jobs since last read from stream
		val := redisClient.XRead(readArgs)
		if val == nil {
			time.Sleep(100 * time.Millisecond)
			continue
		} else{
			streamValues := val.Val()[0].Messages
			lastIDRead = streamValues[len(streamValues) - 1].ID
			dbwriter.Write(streamValues)
			fmt.Println(streamValues)
		}
	}
}

func NewListner(streamClient redis.Client) Listner {
	return Listner{StremClient: streamClient}
}
