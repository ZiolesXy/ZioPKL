package websocket

import (
	"context"
	"log"
	"sync"

	"voca-store/internal/domain/models"
	"voca-store/internal/domain/repository"
)

type Hub struct {
	clients    map[string]map[*Client]bool
	broadcast  chan WSMessage
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
	repo       repository.ChatRepository
}

func NewHub(repo repository.ChatRepository) *Hub {
	return &Hub{
		clients:    make(map[string]map[*Client]bool),
		broadcast:  make(chan WSMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		repo:       repo,
	}
}

// === EXPORTED METHODS (bisa diakses dari package lain) ===

func (h *Hub) RegisterClient(client *Client) {
	h.register <- client
}

func (h *Hub) UnregisterClient(client *Client) {
	h.unregister <- client
}

func (h *Hub) BroadcastToSession(sessionUID string, message WSMessage) {
	h.broadcast <- message
}

func (h *Hub) NotifyAdminsNewRequest(session models.ChatSession) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	notification := NewWSMessage(WSMsgTypeNotification, NotificationPayload{
		Title: "New Chat Request",
		Body:  session.User.Name + " wants to chat",
		Data: map[string]interface{}{
			"type":         "new_chat_request",
			"session_uid":  session.UID,
			"user_id":      session.UserID,
			"user_name":    session.User.Name,
			"user_avatar":  session.User.ProfileImageURL,
			"created_at":   session.CreatedAt,
		},
	})

	if adminClients, ok := h.clients["admin_broadcast"]; ok {
		for client := range adminClients {
			if client.GetRole() == "Admin" {
				select {
				case client.send <- notification:
				default:
				}
			}
		}
	}
}

func (h *Hub) NotifySessionAccepted(sessionUID string, adminName string) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	notification := NewWSMessage(WSMsgTypeHandshakeAck, map[string]interface{}{
		"session_uid": sessionUID,
		"status":      "accepted",
		"admin_name":  adminName,
		"message":     "Admin has accepted your chat request",
	})

	if clients, ok := h.clients[sessionUID]; ok {
		for client := range clients {
			if client.GetRole() != "Admin" {
				select {
				case client.send <- notification:
				default:
				}
			}
		}
	}
}

func (h *Hub) RegisterAdmin(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.clients["admin_broadcast"] == nil {
		h.clients["admin_broadcast"] = make(map[*Client]bool)
	}
	h.clients["admin_broadcast"][client] = true
}

func (h *Hub) UnregisterAdmin(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if clients, ok := h.clients["admin_broadcast"]; ok {
		delete(clients, client)
	}
}

// === INTERNAL RUN LOOP (tidak diubah) ===

func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("WebSocket hub shutting down")
			return

		case client := <-h.register:
			h.mu.Lock()
			if h.clients[client.GetSessionUID()] == nil {
				h.clients[client.GetSessionUID()] = make(map[*Client]bool)
			}
			h.clients[client.GetSessionUID()][client] = true
			h.mu.Unlock()
			log.Printf("Client registered: user %d, session %s", client.GetUserID(), client.GetSessionUID())

		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.clients[client.GetSessionUID()]; ok {
				delete(clients, client)
				if len(clients) == 0 {
					delete(h.clients, client.GetSessionUID())
					log.Printf("Session room closed: %s", client.GetSessionUID())
				}
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.handleBroadcast(message)
		}
	}
}

func (h *Hub) handleBroadcast(message WSMessage) {
	payload, ok := message.Payload.(MessagePayload)
	if !ok {
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	if clients, ok := h.clients[payload.SessionUID]; ok {
		for client := range clients {
			if client.GetUserID() != payload.SenderID {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(clients, client)
				}
			}
		}
	}
}

func (h *Hub) HandleIncomingMessage(client *Client, msg WSMessage) {
	switch msg.Type {
	case WSMsgTypeMessage:
		h.handleChatMessage(client, msg)
	case WSMsgTypeTyping:
		h.handleTyping(client, msg)
	default:
		client.sendError("unknown_type", "Unknown message type")
	}
}

func (h *Hub) handleChatMessage(client *Client, msg WSMessage) {
	payload, ok := msg.Payload.(map[string]interface{})
	if !ok {
		client.sendError("invalid_payload", "Message payload must be an object")
		return
	}

	content, _ := payload["content"].(string)
	msgType, _ := payload["message_type"].(string)
	mediaURL, _ := payload["media_url"].(string)

	if content == "" {
		client.sendError("empty_content", "Message content cannot be empty")
		return
	}

	ctx := context.Background()
	session, err := h.repo.GetSessionByUID(ctx, client.GetSessionUID())
	if err != nil {
		client.sendError("session_not_found", "Chat session not found")
		return
	}

	if session.Status != models.ChatSessionActive {
		client.sendError("session_inactive", "Cannot send message to inactive session")
		return
	}

	// Auto-assign admin current session if not assigned yet
	if client.GetRole() == "Admin" && session.AdminID == nil {
		adminID := client.GetUserID()
		if err := h.repo.UpdateSessionStatus(ctx, session.ID, models.ChatSessionActive, &adminID); err != nil {
			log.Printf("Failed to auto-assign admin to session: %v", err)
		}
	}

	chatMsg := &models.ChatMessage{
		SessionID:   session.ID,
		SenderID:    client.GetUserID(),
		Content:     content,
		MessageType: models.MessageType(msgType),
	}
	if mediaURL != "" {
		chatMsg.MediaURL = &mediaURL
	}

	if err := h.repo.CreateMessage(ctx, chatMsg); err != nil {
		log.Printf("Failed to save message: %v", err)
		client.sendError("save_failed", "Failed to save message")
		return
	}

	if err := h.repo.UpdateLastMessage(ctx, session.ID, content); err != nil {
		log.Printf("Failed to update last message: %v", err)
	}

	senderRole := "user"
	if client.GetRole() == "Admin" {
		senderRole = "admin"
	}

	responseMsg := NewWSMessage(WSMsgTypeMessage, MessagePayload{
		SessionUID:  session.UID,
		SenderID:    client.GetUserID(),
		SenderName:  session.User.Name,
		SenderRole: senderRole,
		Content:     content,
		MessageType: msgType,
		MediaURL:    &mediaURL,
		CreatedAt:   chatMsg.CreatedAt,
	})

	h.mu.RLock()
	if clients, ok := h.clients[client.GetSessionUID()]; ok {
		for c := range clients {
			if c.GetUserID() != client.GetUserID() {
				select {
				case c.send <- responseMsg:
				default:
					close(c.send)
				}
			}
		}
	}
	h.mu.RUnlock()

	client.Send(NewWSMessage(WSMsgTypeMessageAck, map[string]string{
		"message_uid": chatMsg.UID,
		"status":      "sent",
	}))
}

func (h *Hub) handleTyping(client *Client, msg WSMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	typingMsg := NewWSMessage(WSMsgTypeTyping, map[string]interface{}{
		"session_uid": client.GetSessionUID(),
		"user_id":     client.GetUserID(),
		"is_typing":   msg.Payload,
	})

	if clients, ok := h.clients[client.GetSessionUID()]; ok {
		for c := range clients {
			if c.GetUserID() != client.GetUserID() {
				select {
				case c.send <- typingMsg:
				default:
				}
			}
		}
	}
}