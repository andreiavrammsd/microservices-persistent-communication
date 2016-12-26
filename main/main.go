package main

import (
	"net/http"
	"log"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

var servicesQueue *Queue

func main() {
	servicesQueue = NewQueue(config.QueueName)

	go servicesQueue.Consume(func(message Message) {
		for d := range message.Messages {
			service, _ := NewService(d.Body)
			
			if result := service.Call(); result == true {
				d.Ack(false)
			} else {
				d.Nack(false, true)
			}
		}
	})
	
	router := NewRouter()
	log.Fatal(http.ListenAndServe(config.ServerAddress, router))
}
