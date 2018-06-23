package main

import (
	"log"
	"time"

	"github.com/andreiavrammsd/go-rabbitmq"
)

func consumerCallback(d *rabbitmq.Delivery) {
	service, _ := NewService(d.Body)

	if len(service.URL) == 0 {
		log.Print("Ignored: No url provided.")
		if err := d.Ack(false); err != nil {
			log.Println(err)
		}
		return
	}

	err := service.Call()
	if err == nil {
		log.Printf("Success: %s %s.", service.Method, service.URL)
		if err := d.Ack(false); err != nil {
			log.Println(err)
		}
		return
	}

	if service.Requeue {
		log.Printf(
			"Failed: %s %s (%s). Requeuing after %d milliseconds.",
			service.Method,
			service.URL,
			err,
			config.RequeueFailedAfterMilliseconds,
		)
		time.Sleep(time.Millisecond * config.RequeueFailedAfterMilliseconds)
		if err := d.Nack(false, true); err != nil {
			log.Println(err)
		}
		return
	}

	log.Printf(
		"Failed: %s %s. Not requeuing as requested.",
		service.Method,
		service.URL,
	)
	if err := d.Ack(false); err != nil {
		log.Println(err)
	}
}

func ConsumeQueue(numberOfConsumers int) {
	c := &rabbitmq.ConsumerConfig{
		Callback: consumerCallback,
	}
	for i := 1; i <= numberOfConsumers; i++ {
		go func() {
			if err := consumer(c); err != nil {
				log.Println(err)
			}
		}()
	}
}

func consumer(c *rabbitmq.ConsumerConfig) error {
	ch, err := serviceQueueConnection.Channel()
	if err != nil {
		return err
	}

	q, err := ch.Queue(config.ServiceQueueName)
	if err != nil {
		return err
	}

	return q.Consume(c)
}
