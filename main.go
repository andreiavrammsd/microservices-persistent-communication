package main

import (
	"net/http"
	"log"
	"gopkg.in/go-playground/validator.v9"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

var (
	servicesQueue *Queue
	validate *validator.Validate
	config *Config
)

func main() {
	config = NewConfig()

	SetupLogger()
	log.Printf("HTTP server address: %s", config.Server.Address)
	log.Printf("Number of consumers: %d", config.NumberOfConsumers)
	log.Printf("Queue name: %s", config.Queue.Name)
	log.Printf("Log to file: %v", config.FileLogEnabled)
	log.Printf("Fast publish: %v", config.FastPublish)
	log.Printf("Using authorization: %v", len(config.AuthorizationKey) > 0)

	servicesQueue = NewQueue(config.Queue)
	validate = NewValidate()

	consume(config.NumberOfConsumers)

	router := NewRouter()
	log.Fatal(http.ListenAndServe(config.Server.Address, router))
}
