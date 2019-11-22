package dbwriter

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func Write(events []redis.XMessage) {
	clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	eventsCollection := client.Database("sauron_reporters").Collection("event")
	for eventPosition := 0; eventPosition < len(events); eventPosition++ {
		res, err := eventsCollection.InsertOne(context.TODO(), events[eventPosition])
		fmt.Println(res)
		if err != nil {
			fmt.Println(err)
		}
	}
}
