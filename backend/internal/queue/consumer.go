package queue

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type Consumer struct {
	channel *amqp.Channel
}

func NewConsumer(rmq *RabbitMQ) *Consumer {
	return &Consumer{
		channel: rmq.Channel(),
	}
}

type MessageHandler func([]byte) error

func (c *Consumer) Consume(queueName string, handler MessageHandler) error {
	msgs, err := c.channel.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	go func() {
		for msg := range msgs {
			if err := handler(msg.Body); err != nil {
				log.Printf("Error processing message: %v", err)
				msg.Nack(false, true)
			} else {
				msg.Ack(false)
			}
		}
	}()

	log.Printf("Started consuming from queue: %s", queueName)
	return nil
}

func (c *Consumer) ConsumeJSON(queueName string, handler func(interface{}) error) error {
	return c.Consume(queueName, func(body []byte) error {
		var data interface{}
		if err := json.Unmarshal(body, &data); err != nil {
			return fmt.Errorf("failed to unmarshal message: %w", err)
		}
		return handler(data)
	})
}