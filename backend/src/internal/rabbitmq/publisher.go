package rabbitmq

import (
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

// PublishMessage publishes a message to a specified queue
func (r *RabbitMQ) PublishMessage(queueName string, message interface{}) error {
	// Declare the queue (idempotent operation)
	_, err := r.channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue %s: %w", queueName, err)
	}

	// Marshal the message to JSON
	body, err := sonic.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Publish the message
	err = r.channel.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	zap.L().Info("Successfully published message to queue", zap.String("queue", queueName))
	return nil
}
