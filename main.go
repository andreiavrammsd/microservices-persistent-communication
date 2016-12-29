package main

import (
	"log"
	"gopkg.in/go-playground/validator.v9"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var (
	config *Config
	serviceQueueConnection *RabbitMqConnection
	serviceQueue *RabbitMqQueue
	validate *validator.Validate
)

func main() {
	config = GetConfig()

	SetupLogger()

	validate = NewValidate(config.Validation)

	conn, err := RabbitMqNewConnection(config.RabbitMqConfig)
	checkError(err)

	ch, err := conn.GetChannel()
	checkError(err)

	q, err := ch.GetQueue(config.ServiceQueueName)
	checkError(err)

	serviceQueueConnection = conn
	serviceQueue = q

	serverAddress := config.Server.Address
	if (config.Server.Tls) {
		serverAddress = config.Server.AddressTls
	}

	log.Printf("HTTP server address: %s", serverAddress)
	log.Printf("TLS: %v", config.Server.Tls)
	log.Printf("Number of consumers: %d", config.NumberOfConsumers)
	log.Printf("Queue name: %s", config.ServiceQueueName)
	log.Printf("Log to file: %v", config.FileLogEnabled)
	log.Printf("Fast publish: %v", config.FastPublish)
	log.Printf("Using authorization: %v", len(config.AuthorizationKey) > 0)

	ConsumeQueue(config.NumberOfConsumers)

	router := NewRouter()
	StartServer(config.Server, router)
}
