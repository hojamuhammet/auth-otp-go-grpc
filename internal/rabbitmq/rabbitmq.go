package rabbitmq

import (
	"context"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQService struct {
	conn *amqp091.Connection
}

func InitRabbitMQConnection(rabbitMQURL string) (*RabbitMQService, error) {
	conn, err := amqp091.Dial(rabbitMQURL)
	if err != nil {
		return nil, err
	}

	return &RabbitMQService{conn: conn}, nil
}

func (r *RabbitMQService) Close() {
	if r.conn != nil {
		r.conn.Close()
		log.Println("Closed RabbitMQ connection")
	}
}

func (r *RabbitMQService) PublishMessage(ctx context.Context, queueName string, message []byte) error {
	ch, err := r.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

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

	err = ch.PublishWithContext(
		ctx,
		"",        // Exchange (empty string for default exchange)
		queueName, // Routing key (queue name)
		false,     // Mandatory
		false,     // Immediate
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
