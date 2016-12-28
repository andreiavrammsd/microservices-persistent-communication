package main

import (
	"net/http"
	"log"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

var (
	config = NewConfig()
	servicesQueue = NewQueue(config.Queue)
	validate = NewValidate(config.Validation)
)

func main() {
	SetupLogger()

	log.Printf("HTTP server address: %s", config.Server.Address)
	log.Printf("Number of consumers: %d", config.NumberOfConsumers)
	log.Printf("Queue name: %s", config.Queue.Name)
	log.Printf("Log to file: %v", config.FileLogEnabled)
	log.Printf("Fast publish: %v", config.FastPublish)
	log.Printf("Using authorization: %v", len(config.AuthorizationKey) > 0)

	consume(config.NumberOfConsumers)

	router := NewRouter()
	log.Fatal(http.ListenAndServe(config.Server.Address, router))
}
