package main

import (
	"log"

	"github.com/andreiavrammsd/go-rabbitmq"
	"gopkg.in/go-playground/validator.v9"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var (
	config                 *Config
	serviceQueueConnection *rabbitmq.Connection
	serviceQueue           *rabbitmq.Queue
	validate               *validator.Validate
)

func main() {
	config = GetConfig()

	SetupLogger()

	var err error
	validate, err = NewValidate()
	if err != nil {
		log.Fatalf("Error at validation init (%s)", err)
	}

	conn, err := rabbitmq.New(&config.RabbitMqConfig)
	checkError(err)

	ch, err := conn.Channel()
	checkError(err)

	q, err := ch.Queue(config.ServiceQueueName)
	checkError(err)

	serviceQueueConnection = conn
	serviceQueue = q

	serverAddress := config.Server.Address
	if config.Server.TLS {
		serverAddress = config.Server.AddressTLS
	}

	log.Printf("HTTP server address: %s", serverAddress)
	log.Printf("TLS: %v", config.Server.TLS)
	log.Printf("Number of consumers: %d", config.NumberOfConsumers)
	log.Printf("Queue name: %s", config.ServiceQueueName)
	log.Printf("Log to file: %v", config.FileLogEnabled)
	log.Printf("Fast publish: %v", config.FastPublish)
	log.Printf("Using authorization: %v", len(config.AuthorizationKey) > 0)

	ConsumeQueue(config.NumberOfConsumers)

	router := NewRouter()
	StartServer(config.Server, router)
}
