package eventlistener

import (
	"fmt"
	"github.com/step/angmar/pkg/hashClient"
	"log"
	"time"

	"github.com/go-redis/redis"
	"github.com/step/sauron_reporters/pkg/notifier"
	"github.com/step/sauron_reporters/pkg/writer"
	sClient "github.com/step/uruk/pkg/streamClient"
)

type Listener struct {
	SClient    sClient.StreamClient
	HashClient hashClient.HashClient
	Writer     writer.Writer
	Notifier   notifier.Notifier
	Logger     *log.Logger
}

func (l Listener) Start(sName string, r chan<- bool, stop <-chan bool) {
	l.logStart(sName)
	shouldStop := false
	lastReadID := l.HashClient.Get("IDHash", "lastReadID")
	go func() {
		shouldStop = <-stop
	}()

	for {
		// Keep running till asked to stop
		if shouldStop {
			break
		}
		// Take all jobs since last read from stream
		streamValues := l.SClient.Read([]string{sName, lastReadID})
		if len(streamValues) == 0 {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		lastReadID = streamValues[len(streamValues)-1].ID
		l.HashClient.Set("IDHash", "lastReadID", lastReadID)
		l.logRead(lastReadID)

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

func NewListner(sClient sClient.StreamClient, hashClient hashClient.HashClient, w writer.Writer, notifier notifier.Notifier, logger *log.Logger) Listener {
	return Listener{
		SClient:    sClient,
		HashClient: hashClient,
		Writer:     w,
		Notifier:   notifier,
		Logger:     logger,
	}
}
