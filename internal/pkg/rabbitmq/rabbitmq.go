// rabbitmq.go

package rabbitmq

import (
	"log"

	"github.com/streadway/amqp"
)

// RabbitMQService represents the RabbitMQ messaging service.
type RabbitMQService struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

// NewRabbitMQService creates a new instance of the RabbitMQService.
func NewRabbitMQService() (*RabbitMQService, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/") // Update with your RabbitMQ connection details
	if err != nil {
		log.Printf("Failed to connect to RabbitMQ: %v", err)
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		log.Printf("Failed to open a channel: %v", err)
		return nil, err
	}

	// Declare a queue to send phone numbers
	queue, err := channel.QueueDeclare(
		"phone_numbers", // Queue name
		false,           // Durable
		false,           // Delete when unused
		false,           // Exclusive
		false,           // No-wait
		nil,             // Arguments
	)
	if err != nil {
		log.Printf("Failed to declare a queue: %v", err)
		return nil, err
	}

	return &RabbitMQService{conn, channel, queue}, nil
}

// SendPhoneNumber sends a phone number to the RabbitMQ queue.
func (s *RabbitMQService) SendPhoneNumber(phoneNumber string) error {
	err := s.channel.Publish(
		"",                 // Exchange
		s.queue.Name,       // Routing key
		false,              // Mandatory
		false,              // Immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(phoneNumber),
		})
	if err != nil {
		log.Printf("Failed to publish a message to RabbitMQ: %v", err)
		return err
	}

	log.Printf("Sent phone number %s to RabbitMQ", phoneNumber)
	return nil
}

// Close closes the RabbitMQ connection and channel.
func (s *RabbitMQService) Close() {
	if s.channel != nil {
		s.channel.Close()
	}
	if s.conn != nil {
		s.conn.Close()
	}
}
