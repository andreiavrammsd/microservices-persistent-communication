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
	Tls                          bool
	Queue                        QueueConfig
	NumberOfConsumers            int
	RequeueFailedAfterMilliseconds time.Duration
	FileLogEnabled               bool
	LogFile                      string
	FastPublish                  bool
	Validation                   ValidationConfig
}

type ServerConfig struct {
	Tls           bool
	RedirectToTls bool
	Address       string
	AddressTls    string
	CertFile      string
	KeyFile       string
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
	tls, _ := strconv.ParseBool(os.Getenv("TLS"))
	redirectToTls, _ := strconv.ParseBool(os.Getenv("REDIRECT_TO_TLS"))

	numberOfConsumers, err := strconv.Atoi(os.Getenv("QUEUE_NUMBER_OF_CONSUMERS"))
	if numberOfConsumers == 0 || err != nil {
		numberOfConsumers = 3
	}

	requeueFailedAfterMilliseconds, err := strconv.Atoi(os.Getenv("REQUEUE_FAILED_AFTER_MILLISECONDS"))
	if requeueFailedAfterMilliseconds == 0 || err != nil {
		requeueFailedAfterMilliseconds = 5000
	}

	fileLogEnabled, _ := strconv.ParseBool(os.Getenv("FILE_LOG_ENABLED"))
	justPublish, _ := strconv.ParseBool(os.Getenv("JUST_PUBLISH"))

	return &Config{
		Server: ServerConfig{
			Tls: tls,
			RedirectToTls: redirectToTls,
			Address: ":8008",
			AddressTls: ":8009",
			CertFile : "./ssl/server.crt",
			KeyFile : "./ssl/server.key",
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
		RequeueFailedAfterMilliseconds: time.Duration(requeueFailedAfterMilliseconds),
		FileLogEnabled: fileLogEnabled,
		LogFile: "/var/log/microservices-persistent-communication/app.log",
		FastPublish: justPublish,
		Validation: ValidationConfig{
			Protocols: []string{"http", "https"},
			Methods: []string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "CONNECT", "TRACE"},
		},
	}
}
