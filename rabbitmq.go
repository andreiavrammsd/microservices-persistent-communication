package main

import (
	"github.com/streadway/amqp"
	"fmt"
)

type RabbitMqConfig struct {
	Address  string
	Username string
	Password string
}

type RabbitMqConnection struct {
	Config     RabbitMqConfig
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Queue      amqp.Queue
}

type RabbitMqChannel struct {
	Channel *amqp.Channel
}

type RabbitMqQueueConfig struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
}

type RabbitMqQueue struct {
	Channel *amqp.Channel
	Queue   amqp.Queue
}

type RabbitMqMessage struct {
	Body         []byte
	ContentType  string
	DeliveryMode uint8
	Exchange     string
	RoutingKey   string
	Mandatory    bool
	Immediate    bool
}

type RabbitMqDelivery struct {
	Messages <-chan amqp.Delivery
	Queue    amqp.Queue
}

type RabbitMqConsumerConfig struct {
	QueueName string
	Consumer  string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Args      amqp.Table
}

type Consumer func(delivery RabbitMqDelivery)

func RabbitMqNewConnection(config RabbitMqConfig) (*RabbitMqConnection, error) {
	q := &RabbitMqConnection{
		Config: config,
	}

	url := fmt.Sprintf("amqp://%s:%s@%s/", config.Username, config.Password, config.Address)
	conn, err := amqp.Dial(url)

	q.Connection = conn

	return q, err
}

func (q *RabbitMqConnection) GetChannel() (*RabbitMqChannel, error) {
	ch, err := q.Connection.Channel()
	c := &RabbitMqChannel{
		Channel: ch,
	}

	return c, err
}

func (ch *RabbitMqChannel) GetQueue(name string) (*RabbitMqQueue, error) {
	config := &RabbitMqQueueConfig{
		Name: name,
		Durable: true,
		AutoDelete: false,
		Exclusive: false,
		NoWait: false,
		Args : nil,
	}
	queue, err := ch.Channel.QueueDeclare(
		config.Name,
		config.Durable,
		config.AutoDelete,
		config.Exclusive,
		config.NoWait,
		config.Args,
	)
	q := &RabbitMqQueue{
		Channel: ch.Channel,
		Queue: queue,
	}

	return q, err
}

func RabbitMqNewMessage() RabbitMqMessage {
	return RabbitMqMessage{
		ContentType: "text/plain",
		DeliveryMode: amqp.Persistent,
		Exchange: "",
		RoutingKey: "",
		Mandatory: false,
		Immediate: false,
	}
}

func (q *RabbitMqQueue) Publish(m RabbitMqMessage) error {
	err := q.Channel.Publish(
		m.Exchange,
		q.Queue.Name,
		m.Mandatory,
		m.Immediate,
		amqp.Publishing{
			DeliveryMode: m.DeliveryMode,
			ContentType: m.ContentType,
			Body: m.Body,
		},
	)

	return err
}

func (q *RabbitMqQueue) Consume(consumer Consumer) (error) {
	config := &RabbitMqConsumerConfig{
		QueueName: q.Queue.Name,
		Consumer: "",
		AutoAck : false,
		Exclusive: false,
		NoLocal: false,
		NoWait: false,
		Args : nil,
	}
	messages, err := q.Channel.Consume(
		config.QueueName,
		config.Consumer,
		config.AutoAck,
		config.Exclusive,
		config.NoLocal,
		config.NoWait,
		config.Args,
	)

	if err != nil {
		return err
	}

	delivery := RabbitMqDelivery{
		Messages: messages,
		Queue: q.Queue,
	}
	consumer(delivery)

	return nil
}
