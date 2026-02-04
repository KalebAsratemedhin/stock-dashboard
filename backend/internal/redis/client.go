package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	rdb *redis.Client
	ctx context.Context
}

var (
	StockChannel      = "stock_quotes"
	SalesChannel      = "sales"
	UserEventsChannel = "user_events"
	FinancialChannel  = "financial_metrics"
)

func NewClient() (*Client, error) {
	url := getEnv("REDIS_URL", "redis://localhost:6379")
	
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	rdb := redis.NewClient(opt)
	ctx := context.Background()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Println("Redis connected")
	return &Client{rdb: rdb, ctx: ctx}, nil
}

func (c *Client) Publish(channel string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	err = c.rdb.Publish(c.ctx, channel, jsonData).Err()
	if err != nil {
		return fmt.Errorf("failed to publish to Redis: %w", err)
	}

	return nil
}

func (c *Client) Subscribe(channels ...string) *redis.PubSub {
	return c.rdb.Subscribe(c.ctx, channels...)
}

func (c *Client) Close() error {
	return c.rdb.Close()
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

