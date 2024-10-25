package rabbitmq

import (
	"fmt"

	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// NewClient creates a new RabbitMQ client
func NewClient(host, port, username, password, vhost string) (*RabbitMQ, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/%s",
		username, password, host, port, vhost)

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("error connecting to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("error creating RabbitMQ channel: %w", err)
	}

	zap.L().Info("Successfully connected to RabbitMQ")
	return &RabbitMQ{
		conn:    conn,
		channel: ch,
	}, nil
}

// Close closes the underlying RabbitMQ connection and channel
func (r *RabbitMQ) Close() error {
	if r.channel != nil {
		if err := r.channel.Close(); err != nil {
			return fmt.Errorf("error closing channel: %w", err)
		}
	}
	if r.conn != nil {
		if err := r.conn.Close(); err != nil {
			return fmt.Errorf("error closing connection: %w", err)
		}
	}
	return nil
}

// GetChannel returns the underlying RabbitMQ channel
func (r *RabbitMQ) GetChannel() *amqp.Channel {
	return r.channel
}

// GetConnection returns the underlying RabbitMQ connection
func (r *RabbitMQ) GetConnection() *amqp.Connection {
	return r.conn
}