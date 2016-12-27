package main

import "time"

type Config struct {
	ServerAddress                string
	Queue                        QueueConfig
	NumberOfConsumers            int
	RetryFailedAfterMilliseconds time.Duration
}

type QueueConfig struct {
	Address  string
	Username string
	Password string
	Name     string
}

var config = &Config{
	ServerAddress: ":8008",
	Queue: QueueConfig{
		Address: "127.0.0.1:5672",
		Username: "WhLSCKgkzL66aAvQ",
		Password: "Ayxae5yNGUtQVSufkp44xPgTJpaBeQKS",
		Name: "services",
	},
	NumberOfConsumers: 3,
	RetryFailedAfterMilliseconds: 5000,
}
