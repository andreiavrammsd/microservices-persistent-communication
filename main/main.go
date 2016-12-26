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

func main() {
	go NewQueue(config.QueueName).Consume(func (message Message) {
		for d := range message.Messages {
			service, _ := NewService(d.Body)
			service.Call()
		}
	})

	router := NewRouter()
	log.Fatal(http.ListenAndServe(config.ServerAddress, router))
}
