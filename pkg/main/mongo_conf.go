package main

import (
	"context"
	"flag"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"

	"github.com/step/sauron_reporters/pkg/writer"
)

var db string
var table string

func init() {
	flag.StringVar(&db, "db", "sauron_reporters", "Mongo `database` to transact with")
	flag.StringVar(&table, "table", "events", "`table name` to fetch data")
}

func getMongoWriter(logger *log.Logger) writer.MongoDbWriter {
	userName := os.Getenv("MONGO_USER_NAME")
	password := os.Getenv("MONGO_PASSWORD")
	mongoURI := fmt.Sprintf("mongodb+srv://%s:%s@cluster0-nvzs5.mongodb.net/test?retryWrites=true&w=majority", userName, password)

	client, _ := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err := client.Connect(ctx)

	if err != nil {
		fmt.Println("Unable to connect due to => ", err)
	}

	return writer.MongoDbWriter{
		DB:     db,
		Table:  table,
		Client: client,
		Logger: logger,
	}
}
