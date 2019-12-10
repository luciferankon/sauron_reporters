package eventlistener

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/step/sauron_reporters/pkg/notifier"
	"github.com/step/sauron_reporters/pkg/writer"
	sClient "github.com/step/uruk/pkg/streamClient"
)

type Listner struct {
	SClient  sClient.StreamClient
	Writer   writer.Writer
	Notifier notifier.Notifier
}

func (l Listner) Start(streamName string, r chan<- bool, stop <-chan bool) {
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
		// Take all jobs since last read from stream
		streamValues := l.SClient.Read([]string{streamName, lastIDRead})
		if len(streamValues) == 0 {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		lastIDRead = streamValues[len(streamValues)-1].ID

		for _, value := range streamValues {
			if value.Values["type"] == "job_complete" {
				event := redis.XMessage{
					ID:     value.ID,
					Values: value.Values,
				}
				l.Writer.Write(event.Values)
				l.Notifier.Notify(event.Values)
			}
		}
	}
}

func NewListner(sClient sClient.StreamClient, w writer.Writer, notifier notifier.Notifier) Listner {
	return Listner{
		SClient: sClient,
		Writer:  w,
		Notifier: notifier,
	}
}
