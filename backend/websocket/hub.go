package websocket

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

type Client struct {
	Conn   *websocket.Conn
	PollID string
}

type Hub struct {
	clients map[string]map[*Client]bool
	mu      sync.RWMutex
	redis   *redis.Client
}

func NewHub(redisClient *redis.Client) *Hub {
	return &Hub{
		clients: make(map[string]map[*Client]bool),
		redis:   redisClient,
	}
}

func (h *Hub) AddClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.clients[client.PollID] == nil {
		h.clients[client.PollID] = make(map[*Client]bool)
	}

	h.clients[client.PollID][client] = true
}

func (h *Hub) RemoveClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if pollClients, ok := h.clients[client.PollID]; ok {
		delete(pollClients, client)

		if len(pollClients) == 0 {
			delete(h.clients, client.PollID)
		}
	}
}

func (h *Hub) Broadcast(pollID string, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients[pollID] {
		if err := client.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("websocket write error: %v", err)
		}
	}
}

// StartRedisSubscriber listens to every poll update published by Redis.
func (h *Hub) StartRedisSubscriber(ctx context.Context) {
	pubsub := h.redis.PSubscribe(ctx, "poll:*:updates")
	defer pubsub.Close()

	log.Println("Redis Pub/Sub subscriber started")

	for {
		message, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}

			log.Printf("Redis Pub/Sub error: %v", err)
			continue
		}

		// Channel format:
		// poll:{pollID}:updates
		parts := strings.Split(message.Channel, ":")

		if len(parts) != 3 {
			continue
		}

		pollID := parts[1]

		// Make sure the published data is valid JSON.
		var payload json.RawMessage

		if err := json.Unmarshal([]byte(message.Payload), &payload); err != nil {
			log.Printf("invalid Redis message: %v", err)
			continue
		}

		h.Broadcast(pollID, []byte(message.Payload))
	}
}
