package main

import (
	"flag"
	"log"

	"github.com/step/sauron_reporters/pkg/writer"
)

var mongoURI string
var db string
var table string

func init() {
	flag.StringVar(&mongoURI, "mongo-uri", "mongodb://127.0.0.1:27017", "`address` of Mongo host to connect to")
	flag.StringVar(&db, "db", "sauron_reporters", "Mongo `database` to transact with")
	flag.StringVar(&table, "table", "events", "`table name` to fetch data")
}

func getMongoWriter(logger *log.Logger) writer.MongoDbWriter {
	return writer.MongoDbWriter{
		URI:    mongoURI,
		DB:     db,
		Table:  table,
		Logger: logger,
	}
}
