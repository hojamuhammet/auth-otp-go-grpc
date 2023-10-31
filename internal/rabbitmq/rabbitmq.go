package rabbitmq

import (
	"context"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

// RabbitMQService is a service for interacting with RabbitMQ.
type RabbitMQService struct {
    conn *amqp091.Connection
}

// NewRabbitMQService initializes a RabbitMQ service.
func InitRabbitMQConnection(rabbitMQURL string) (*RabbitMQService, error) {
    // Establish a connection to RabbitMQ
    conn, err := amqp091.Dial(rabbitMQURL)
    if err != nil {
        return nil, err
    }

    log.Println("Connected to RabbitMQ")
    return &RabbitMQService{conn: conn}, nil
}

// Close closes the RabbitMQ connection.
func (r *RabbitMQService) Close() {
    if r.conn != nil {
        r.conn.Close()
        log.Println("Closed RabbitMQ connection")
    }
}

// PublishMessage publishes a message to the specified queue.
func (r *RabbitMQService) PublishMessage(ctx context.Context, queueName string, message []byte) error {
    // Create a channel for the connection
    ch, err := r.conn.Channel()
    if err != nil {
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
        return err
    }

    // Publish the message to the queue
    err = ch.PublishWithContext(
        ctx,
        "",         // Exchange (empty string for default exchange)
        queueName,   // Routing key (queue name)
        false,      // Mandatory
        false,      // Immediate
        amqp091.Publishing{
            ContentType: "application/grpc+proto",
            Body:        message,
        },
    )
    if err != nil {
        return err
    }

    log.Printf("Published message to queue %s: %s", queueName, message)
    return nil
}
