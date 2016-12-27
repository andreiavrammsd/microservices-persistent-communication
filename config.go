package main

import (
	"time"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Server                       ServerConfig
	Queue                        QueueConfig
	NumberOfConsumers            int
	RetryFailedAfterMilliseconds time.Duration
	LogFile                      string
}

type ServerConfig struct {
	Address string
}

type QueueConfig struct {
	Address  string
	Username string
	Password string
	Name     string
}

func NewConfig() *Config {
	numberOfConsumers, err := strconv.Atoi(os.Getenv("QUEUE_NUMBER_OF_CONSUMERS"))
	if numberOfConsumers == 0 || err != nil {
		numberOfConsumers = 3
	}

	retryFailedAfterMilliseconds, err := strconv.Atoi(os.Getenv("RETRY_FAILED_AFTER_MILLISECONDS"))
	if retryFailedAfterMilliseconds == 0 || err != nil {
		retryFailedAfterMilliseconds = 5000
	}

	return &Config{
		Server: ServerConfig{
			Address: ":8008",
		},
		Queue: QueueConfig{
			Address: fmt.Sprintf("%s:%s", os.Getenv("QUEUE_HOST"), os.Getenv("QUEUE_PORT")),
			Username: os.Getenv("RABBITMQ_DEFAULT_USER"),
			Password: os.Getenv("RABBITMQ_DEFAULT_PASS"),
			Name: os.Getenv("QUEUE_NAME"),
		},
		NumberOfConsumers: numberOfConsumers,
		RetryFailedAfterMilliseconds: time.Duration(retryFailedAfterMilliseconds),
		LogFile: "/var/log/microservices-persistent-communication/app.log",
	}
}

var config = NewConfig()
