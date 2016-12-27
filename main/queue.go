package main

import (
	"github.com/streadway/amqp"
	"log"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

type Queue struct {
	Name string
	Channel *amqp.Channel
	Queue amqp.Queue
}

type Message struct {
	Messages <-chan amqp.Delivery
	Queue amqp.Queue
}

type Consumer func(message Message)

func (queue *Queue) Connect() {
	conn, err := amqp.Dial(config.QueueServerAddress)
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	q, errDeclare := ch.QueueDeclare(
		queue.Name, // name
		true, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil, // arguments
	)
	failOnError(errDeclare, "Failed to declare a queue")
	
	queue.Channel = ch
	queue.Queue = q
}

func (queue *Queue) Publish(message []byte) {
	e := queue.Channel.Publish(
		"", // exchange
		queue.Name, // routing key,
		false, // mandatory,
		false, // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType: "text/plain",
			Body: message,
		},
	)

	failOnError(e, "Failed to publish a message")
	log.Printf("Received: %s", string(message))
}

func (queue *Queue) Consume(consumer Consumer) {
	messages, err := queue.Channel.Consume(
		queue.Name, // queue
		"", // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil, // args
	)
	failOnError(err, "Failed to register a consumer")
	
	message := Message{
		Messages: messages,
		Queue: queue.Queue,
	}
	consumer(message)
}

func NewQueue(queueName string) *Queue {
	queue := &Queue{
		Name: queueName,
	}
	queue.Connect()
	return queue
}
