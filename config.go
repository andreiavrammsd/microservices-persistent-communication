package main

import (
	"os"
	"strconv"
	"time"

	"github.com/andreiavrammsd/go-rabbitmq"
)

type Config struct {
	Server                         ServerConfig
	AuthorizationHeader            string
	AuthorizationKey               string
	TLS                            bool
	RabbitMqConfig                 rabbitmq.Config
	ServiceQueueName               string
	NumberOfConsumers              int
	RequeueFailedAfterMilliseconds time.Duration
	FileLogEnabled                 bool
	LogFile                        string
	FastPublish                    bool
	Validation                     ValidationConfig
}

type ServerConfig struct {
	TLS           bool
	RedirectToTLS bool
	Address       string
	AddressTLS    string
	CertFile      string
	KeyFile       string
}

type ValidationConfig struct {
	Protocols []string
	Methods   []string
}

func GetConfig() *Config {
	tls, _ := strconv.ParseBool(os.Getenv("TLS"))
	redirectToTLS, _ := strconv.ParseBool(os.Getenv("REDIRECT_TO_TLS"))

	queuePort, err := strconv.Atoi(os.Getenv("QUEUE_HOST"))
	if err != nil {
		queuePort = 5672
	}

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
			TLS:           tls,
			RedirectToTLS: redirectToTLS,
			Address:       ":8008",
			AddressTLS:    ":8009",
			CertFile:      "./ssl/server.crt",
			KeyFile:       "./ssl/server.key",
		},
		AuthorizationHeader: os.Getenv("AUTHORIZATION_HEADER"),
		AuthorizationKey:    os.Getenv("AUTHORIZATION_KEY"),
		RabbitMqConfig: rabbitmq.Config{
			Scheme:   os.Getenv("QUEUE_SCHEME"),
			Host:     os.Getenv("QUEUE_HOST"),
			Port:     queuePort,
			Username: os.Getenv("RABBITMQ_DEFAULT_USER"),
			Password: os.Getenv("RABBITMQ_DEFAULT_PASS"),
			Vhost:    os.Getenv("RABBITMQ_DEFAULT_VHOST"),
		},
		ServiceQueueName:               os.Getenv("QUEUE_NAME"),
		NumberOfConsumers:              numberOfConsumers,
		RequeueFailedAfterMilliseconds: time.Duration(requeueFailedAfterMilliseconds),
		FileLogEnabled:                 fileLogEnabled,
		LogFile:                        "/var/log/microservices-persistent-communication/app.log",
		FastPublish:                    justPublish,
		Validation: ValidationConfig{
			Protocols: []string{"http", "https"},
			Methods:   []string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "CONNECT", "TRACE"},
		},
	}
}
