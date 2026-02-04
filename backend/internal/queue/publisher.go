package queue

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

type Publisher struct {
	channel *amqp.Channel
}

func NewPublisher(rmq *RabbitMQ) *Publisher {
	return &Publisher{
		channel: rmq.Channel(),
	}
}

func (p *Publisher) PublishStockQuote(queue string, data interface{}) error {
	return p.publish(queue, data)
}

func (p *Publisher) PublishSale(queue string, data interface{}) error {
	return p.publish(queue, data)
}

func (p *Publisher) PublishUserEvent(queue string, data interface{}) error {
	return p.publish(queue, data)
}

func (p *Publisher) PublishFinancialMetric(queue string, data interface{}) error {
	return p.publish(queue, data)
}

func (p *Publisher) publish(queueName string, data interface{}) error {
	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = p.channel.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf("Published message to queue: %s", queueName)
	return nil
}

func (p *Publisher) PublishBatch(queueName string, items []interface{}) error {
	for _, item := range items {
		if err := p.publish(queueName, item); err != nil {
			return fmt.Errorf("failed to publish batch item: %w", err)
		}
	}
	return nil
}