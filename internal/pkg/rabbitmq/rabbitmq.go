package rabbitmq

import (
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