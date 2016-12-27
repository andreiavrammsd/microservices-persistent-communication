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

var servicesQueue *Queue
var validate *validator.Validate

func main() {
	log.Printf("HTTP server address: %s", config.ServerAddress)
	log.Printf("Number of consumers: %d", config.NumberOfConsumers)
	log.Printf("Queue name: %s", config.Queue.Name)

	servicesQueue = NewQueue(config.Queue)
	validate = NewValidate()

	consume(config.NumberOfConsumers)

	router := NewRouter()
	log.Fatal(http.ListenAndServe(config.ServerAddress, router))
}
