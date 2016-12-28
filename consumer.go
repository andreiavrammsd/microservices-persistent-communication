package main

import (
	"time"
	"log"
)

func consumer() {
	servicesQueue.Consume(func(message Message) {
		for d := range message.Messages {
			service, _ := NewService(d.Body)

			if success := service.Call(); success == true {
				log.Printf("Success: %s.", service.Url)
				d.Ack(false)
			} else {
				if service.Retry {
					log.Printf(
						"Failed: %s. Retrying after %d milliseconds.",
						service.Url,
						config.RetryFailedAfterMilliseconds,
					)
					time.Sleep(time.Millisecond * config.RetryFailedAfterMilliseconds)
					d.Nack(false, true)
				} else {
					log.Printf("Failed: %s. Not retrying.", service.Url)
					d.Ack(false)
				}
			}
		}
	})
}

func consume(numberOfConsumers int) {
	for i := 1; i <= numberOfConsumers; i++ {
		go consumer()
	}
}
