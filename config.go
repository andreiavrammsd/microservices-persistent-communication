package main

import (
	"time"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Server                       ServerConfig
	AuthorizationHeader          string
	AuthorizationKey             string
	Queue                        QueueConfig
	NumberOfConsumers            int
	RetryFailedAfterMilliseconds time.Duration
	FileLogEnabled               bool
	LogFile                      string
	FastPublish                  bool
	Validation                   ValidationConfig
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

type ValidationConfig struct {
	Protocols []string
	Methods   []string
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

	fileLogEnabled, _ := strconv.ParseBool(os.Getenv("FILE_LOG_ENABLED"))

	justPublish, _ := strconv.ParseBool(os.Getenv("JUST_PUBLISH"))

	return &Config{
		Server: ServerConfig{
			Address: ":8008",
		},
		AuthorizationHeader: os.Getenv("AUTHORIZATION_HEADER"),
		AuthorizationKey: os.Getenv("AUTHORIZATION_KEY"),
		Queue: QueueConfig{
			Address: fmt.Sprintf("%s:%s", os.Getenv("QUEUE_HOST"), os.Getenv("QUEUE_PORT")),
			Username: os.Getenv("RABBITMQ_DEFAULT_USER"),
			Password: os.Getenv("RABBITMQ_DEFAULT_PASS"),
			Name: os.Getenv("QUEUE_NAME"),
		},
		NumberOfConsumers: numberOfConsumers,
		RetryFailedAfterMilliseconds: time.Duration(retryFailedAfterMilliseconds),
		FileLogEnabled: fileLogEnabled,
		LogFile: "/var/log/microservices-persistent-communication/app.log",
		FastPublish: justPublish,
		Validation: ValidationConfig{
			Protocols: []string{"http", "https"},
			Methods: []string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "CONNECT", "TRACE"},
		},
	}
}
