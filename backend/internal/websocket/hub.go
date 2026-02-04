package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	subs     map[string]bool
	subsMu   sync.RWMutex
}

type Message struct {
	Type    string      `json:"type"`
	Channel string      `json:"channel,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("Client connected. Total clients: %d", len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
			log.Printf("Client disconnected. Total clients: %d", len(h.clients))

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) Broadcast(channel string, data interface{}) {
	msg := Message{
		Type:    "data",
		Channel: channel,
		Data:    data,
	}
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}
	h.broadcast <- jsonMsg
}

func (c *Client) Subscribe(channel string) {
	c.subsMu.Lock()
	c.subs[channel] = true
	c.subsMu.Unlock()
}

func (c *Client) Unsubscribe(channel string) {
	c.subsMu.Lock()
	delete(c.subs, channel)
	c.subsMu.Unlock()
}

func (c *Client) IsSubscribed(channel string) bool {
	c.subsMu.RLock()
	defer c.subsMu.RUnlock()
	return c.subs[channel]
}

func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			continue
		}

		switch msg.Type {
		case "subscribe":
			c.Subscribe(msg.Channel)
			response := Message{Type: "subscribed", Channel: msg.Channel}
			jsonResp, _ := json.Marshal(response)
			c.send <- jsonResp
		case "unsubscribe":
			c.Unsubscribe(msg.Channel)
			response := Message{Type: "unsubscribed", Channel: msg.Channel}
			jsonResp, _ := json.Marshal(response)
			c.send <- jsonResp
		}
	}
}

func (c *Client) WritePump() {
	defer c.conn.Close()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("Error writing message: %v", err)
				return
			}
		}
	}
}

