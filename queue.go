package main

import (
	"github.com/streadway/amqp"
	"log"
	"fmt"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

type Queue struct {
	Config  QueueConfig
	Channel *amqp.Channel
	Queue   amqp.Queue
}

type Message struct {
	Messages <-chan amqp.Delivery
	Queue    amqp.Queue
}

type Consumer func(message Message)

func (queue *Queue) Connect() {
	url := fmt.Sprintf(
		"amqp://%s:%s@%s/",
		queue.Config.Username,
		queue.Config.Password,
		queue.Config.Address,
	)
	conn, err := amqp.Dial(url)
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	q, errDeclare := ch.QueueDeclare(
		queue.Config.Name, // name
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
		queue.Config.Name, // routing key,
		false, // mandatory,
		false, // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType: "text/plain",
			Body: message,
		},
	)

	failOnError(e, "Failed to publish a message")
}

func (queue *Queue) Consume(consumer Consumer) {
	messages, err := queue.Channel.Consume(
		queue.Config.Name, // queue
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

func NewQueue(config QueueConfig) Queue {
	queue := Queue{
		Config: config,
	}
	queue.Connect()
	return queue
}
