package websocket

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/redis/go-redis/v9"

	"server/internal/domain/dto"
	"server/internal/service"
)

type Hub struct {
	clients    map[uint]map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	chat       *service.ChatService
}

var (
	redisMu     sync.RWMutex
	redisClient *redis.Client
)

func NewHub(chat *service.ChatService) *Hub {
	return &Hub{
		clients:    make(map[uint]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
		chat:       chat,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			if h.clients[client.userID] == nil {
				h.clients[client.userID] = make(map[*Client]bool)
			}
			h.clients[client.userID][client] = true
			updateActiveConnectionsCount(h.connectionCount())
		case client := <-h.unregister:
			if clients, ok := h.clients[client.userID]; ok {
				if _, exists := clients[client]; exists {
					delete(clients, client)
					close(client.send)
				}
				if len(clients) == 0 {
					delete(h.clients, client.userID)
				}
			}
			updateActiveConnectionsCount(h.connectionCount())
		case payload := <-h.broadcast:
			var out dto.WSOutgoingMessage
			if err := json.Unmarshal(payload, &out); err != nil {
				log.Println(err)
				continue
			}
			for _, userID := range []uint{out.SenderID, out.ReceiverID} {
				if userClients, ok := h.clients[userID]; ok {
					for client := range userClients {
						select {
						case client.send <- payload:
						default:
							close(client.send)
							delete(userClients, client)
						}
					}
				}
			}
		}
	}
}

func (h *Hub) connectionCount() int {
	total := 0
	for _, clients := range h.clients {
		total += len(clients)
	}
	return total
}

func (h *Hub) SaveAndBroadcast(senderID uint, receiverID uint, content string) error {
	message, err := h.chat.SaveMessage(senderID, receiverID, content)
	if err != nil {
		return err
	}
	response, err := json.Marshal(dto.WSOutgoingMessage{
		ID:         message.ID,
		SenderID:   message.SenderID,
		ReceiverID: message.ReceiverID,
		Content:    message.Content,
		CreatedAt:  message.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
	if err != nil {
		return err
	}
	h.broadcast <- response
	return nil
}

func updateActiveConnectionsCount(count int) {
	client := sharedRedis()
	if client == nil {
		return
	}
	if err := client.Set(context.Background(), "active_chats", count, 0).Err(); err != nil {
		log.Println(err)
	}
}

func SetRedisClient(client *redis.Client) {
	redisMu.Lock()
	defer redisMu.Unlock()
	redisClient = client
}

func sharedRedis() *redis.Client {
	redisMu.RLock()
	defer redisMu.RUnlock()
	return redisClient
}
