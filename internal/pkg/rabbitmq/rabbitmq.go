package rabbitmq

import (
	"log"

	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQ() (*RabbitMQ, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQ{
		conn:    conn,
		channel: ch,
	}, nil
}

func (rmq *RabbitMQ) Setup() error {
	exchangeName := "otp_exchange"
	err := rmq.channel.ExchangeDeclare(
		exchangeName,
		"direct",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	queueName := "otp_verification_queue"
	_, err = rmq.channel.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	routingKey := "otp_verification"
	err = rmq.channel.QueueBind(
		queueName,
		routingKey,
		exchangeName,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	log.Printf("RabbitMQ setup completed. Exchange: %s, Queue: %s, Routing Key: %s", exchangeName, queueName, routingKey)
	return nil
}

func (rmq *RabbitMQ) ConsumeOTPVerification() (<-chan amqp.Delivery, error) {
	queueName := "otp_verification_queue"
	msgs, err := rmq.channel.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}

func (rmq *RabbitMQ) Close() {
	if rmq.channel != nil {
		rmq.channel.Close()
	}
	if rmq.conn != nil {
		rmq.conn.Close()
	}
}
