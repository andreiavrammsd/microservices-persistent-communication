package main

import (
	"log"
	"time"
)

func ConsumeQueue(numberOfConsumers int) {
	for i := 1; i <= numberOfConsumers; i++ {
		go consumer()
	}
}

func consumer() {
	ch, _ := serviceQueueConnection.GetChannel()
	q, _ := ch.GetQueue(config.ServiceQueueName)

	q.Consume(func(delivery RabbitMqDelivery) {
		for d := range delivery.Messages {
			service, _ := NewService(d.Body)

			if len(service.Url) == 0 {
				log.Print("Ignored: No url provided.")
				d.Ack(false)
			} else {
				if success := service.Call(); success == true {
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
	})
}

