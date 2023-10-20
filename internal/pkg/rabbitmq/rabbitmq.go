package rabbitmq

import (
	"context"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

func InitializeRabbitMQ(rabbitMQURL string) (*amqp091.Connection, error) {
	// Establish a connection to RabbitMQ
	conn, err := amqp091.Dial(rabbitMQURL)
	if err != nil {
		log.Printf("Failed to connect to RabbitMQ: %v", err)
		return nil, err
	}

	log.Println("Connected to RabbitMQ")

	return conn, nil
}

func PublishMessage(ctx context.Context, conn *amqp091.Connection, queueName, message string) error {
	// Create a channel for the connection
	ch, err := conn.Channel()
	if err != nil {
		log.Printf("Failed to open a channel: %v", err)
		return err
	}
	defer ch.Close()

	// Declare a queue to ensure it exists
	_, err = ch.QueueDeclare(
		queueName, // Queue name
		true,      // Durable (messages survive server restart)
		false,     // Delete when unused
		false,     // Exclusive (only one consumer at a time)
		false,     // No-wait
		nil,       // Arguments
	)
	if err != nil {
		log.Printf("Failed to declare a queue: %v", err)
		return err
	}

	// Publish the message to the queue
	err = ch.PublishWithContext(
		ctx,
		"",        // Exchange (empty string for default exchange)
		queueName, // Routing key (queue name)
		false,     // Mandatory
		false,     // Immediate
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		log.Printf("Failed to publish a message: %v", err)
		return err
	}

	log.Printf("Published message to queue %s: %s", queueName, message)

	return nil
}
