package main

import "time"

type Config struct {
	ServerAddress string
	QueueServerAddress string
	QueueName string
	NumberOfConsumers int
	RetryFailedAfterMilliseconds time.Duration
}

var config = Config{
	ServerAddress: ":8008",
	QueueServerAddress: "amqp://WhLSCKgkzL66aAvQ:Ayxae5yNGUtQVSufkp44xPgTJpaBeQKS@127.0.0.1:5672/",
	QueueName: "services",
	NumberOfConsumers: 3,
	RetryFailedAfterMilliseconds: 5000,
}
