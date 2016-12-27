package main

func consumer() {
	servicesQueue.Consume(func(message Message) {
		for d := range message.Messages {
			service, _ := NewService(d.Body)

			if success := service.Call(); success == true {
				d.Ack(false)
			} else {
				d.Nack(false, true)
			}
		}
	})
}

func consume(numberOfConsumers int) {
	for i := 1; i <= numberOfConsumers; i++ {
		go consumer()
	}
}
