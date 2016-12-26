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
}

type Message struct {
	Messages <-chan amqp.Delivery
	Queue amqp.Queue
}
type Consumer func(message Message)

func (queue *Queue) Publish(message []byte) {
	conn, err := amqp.Dial(config.QueueServerAddress)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failled to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queue.Name, // name
		true, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil, // arguments
	)
	failOnError(err, "Failed to declare a queue")

	e := ch.Publish(
		"", // exchange
		q.Name, // routing key,
		false, // mandatory,
		false, // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType: "text/plain",
			Body: message,
		},
	)

	failOnError(e, "Failed to publish a message")
	log.Printf(" [x] Send %s", string(message))
}

func (queue *Queue) Consume(consumer Consumer) {
	conn, err := amqp.Dial(config.QueueServerAddress)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failled to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queue.Name, // name
		true, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil, // arguments
	)
	failOnError(err, "Failed to declare a queue")

	messages, err := ch.Consume(
		q.Name, // queue
		"", // consumer
		true, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil, // args
	)
	failOnError(err, "Failed to register a consumer")
	
	message := Message{
		Messages: messages,
		Queue: q,
	}
	consumer(message)
}

func NewQueue(queueName string) *Queue {
	return &Queue{
		Name: queueName,
	};
}
