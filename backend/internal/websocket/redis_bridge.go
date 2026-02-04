package websocket

import (
	"encoding/json"
	"log"

	"github.com/KalebAsratemedhi/stock-dashboard/backend/internal/redis"
)

type RedisBridge struct {
	hub    *Hub
	client *redis.Client
}

func NewRedisBridge(hub *Hub, client *redis.Client) *RedisBridge {
	return &RedisBridge{
		hub:    hub,
		client: client,
	}
}

func (rb *RedisBridge) Start() {
	pubsub := rb.client.Subscribe(
		redis.StockChannel,
		redis.SalesChannel,
		redis.UserEventsChannel,
		redis.FinancialChannel,
	)
	defer pubsub.Close()

	ch := pubsub.Channel()

	for msg := range ch {
		var data interface{}
		if err := json.Unmarshal([]byte(msg.Payload), &data); err != nil {
			log.Printf("Error unmarshaling Redis message: %v", err)
			continue
		}

		rb.hub.Broadcast(msg.Channel, data)
	}
}

