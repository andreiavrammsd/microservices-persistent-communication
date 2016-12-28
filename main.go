package main

import (
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

	serverAddress := config.Server.Address
	if (config.Server.Tls) {
		serverAddress = config.Server.AddressTls
	}

	log.Printf("HTTP server address: %s", serverAddress)
	log.Printf("TLS: %v", config.Server.Tls)
	log.Printf("Number of consumers: %d", config.NumberOfConsumers)
	log.Printf("Queue name: %s", config.Queue.Name)
	log.Printf("Log to file: %v", config.FileLogEnabled)
	log.Printf("Fast publish: %v", config.FastPublish)
	log.Printf("Using authorization: %v", len(config.AuthorizationKey) > 0)

	ConsumeQueue(config.NumberOfConsumers)

	router := NewRouter()
	StartServer(config.Server, router)
}
