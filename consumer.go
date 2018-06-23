package main

import (
	"log"
	"time"

	"github.com/andreiavrammsd/go-rabbitmq"
)

func consumerCallback(d *rabbitmq.Delivery) {
	service, _ := NewService(d.Body)

	if len(service.Url) == 0 {
		log.Print("Ignored: No url provided.")
		d.Ack(false)
	} else {
		if success := service.Call(); success {
			log.Printf("Success: %s %s.", service.Method, service.Url)
			d.Ack(false)
		} else {
			if service.Requeue {
				log.Printf(
					"Failed: %s %s. Requeuing after %d milliseconds.",
					service.Method,
					service.Url,
					config.RequeueFailedAfterMilliseconds,
				)
				time.Sleep(time.Millisecond * config.RequeueFailedAfterMilliseconds)
				d.Nack(false, true)
			} else {
				log.Printf(
					"Failed: %s %s. Not requeuing as requested.",
					service.Method,
					service.Url,
				)
				d.Ack(false)
			}
		}
	}
}

func ConsumeQueue(numberOfConsumers int) {
	c := &rabbitmq.ConsumerConfig{
		Callback: consumerCallback,
	}
	for i := 1; i <= numberOfConsumers; i++ {
		go consumer(c)
	}
}

func consumer(c *rabbitmq.ConsumerConfig) {
	ch, _ := serviceQueueConnection.Channel()
	q, _ := ch.Queue(config.ServiceQueueName)
	q.Consume(c)
}
