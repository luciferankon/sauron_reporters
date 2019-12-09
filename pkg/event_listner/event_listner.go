package eventlistner

import (
	"time"

	"github.com/go-redis/redis"
	dbwriter "github.com/step/sauron_reporters/pkg/db_writer"
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
		}
		streamValues := val.Val()[0].Messages
		lastIDRead = streamValues[len(streamValues)-1].ID
		var jobCompleteEvents []redis.XMessage
		for _, value := range streamValues {
			if value.Values["type"] == "job_complete" {
				jobCompleteEvents = append(jobCompleteEvents, value)
			}
		}
		dbwriter.Write(jobCompleteEvents)
	}
}

func NewListner(streamClient redis.Client) Listner {
	return Listner{StremClient: streamClient}
}
