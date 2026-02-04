package queue

import (
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

var (
	StockQueue      = "stock_quotes"
	SalesQueue      = "sales"
	UserEventsQueue = "user_events"
	FinancialQueue  = "financial_metrics"
)

func NewRabbitMQ() (*RabbitMQ, error) {
	url := getEnv("RABBITMQ_URL", "amqp://admin:admin@localhost:5672/")
	
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	rmq := &RabbitMQ{
		conn:    conn,
		channel: channel,
	}

	if err := rmq.setupQueues(); err != nil {
		rmq.Close()
		return nil, fmt.Errorf("failed to setup queues: %w", err)
	}

	log.Println("RabbitMQ connected and queues configured")
	return rmq, nil
}

func (r *RabbitMQ) setupQueues() error {
	queues := []string{
		StockQueue,
		SalesQueue,
		UserEventsQueue,
		FinancialQueue,
	}

	for _, queueName := range queues {
		_, err := r.channel.QueueDeclare(
			queueName,
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return fmt.Errorf("failed to declare queue %s: %w", queueName, err)
		}
	}

	return nil
}

func (r *RabbitMQ) Channel() *amqp.Channel {
	return r.channel
}

func (r *RabbitMQ) Close() error {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}

func (r *RabbitMQ) HealthCheck() error {
	if r.conn == nil || r.conn.IsClosed() {
		return fmt.Errorf("RabbitMQ connection is closed")
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}