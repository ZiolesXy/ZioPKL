package websocket

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512000
)

// === EXPORTED UPGRADER (bisa diakses dari package lain) ===
var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // TODO: Restrict to allowed origins in production
	},
}

// === HELPER FUNCTION UNTUK UPGRADE ===
func UpgradeConnection(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	return Upgrader.Upgrade(w, r, nil)
}

// === CLIENT STRUCT (tidak diubah) ===
type Client struct {
	hub        *Hub
	conn       *websocket.Conn
	send       chan WSMessage
	sessionUID string
	userID     uint
	role       string
	mu         sync.RWMutex
	done       chan struct{}
	once       sync.Once
}

func NewClient(hub *Hub, conn *websocket.Conn, sessionUID string, userID uint, role string) *Client {
	return &Client{
		hub:        hub,
		conn:       conn,
		sessionUID: sessionUID,
		userID:     userID,
		role:       role,
		send:       make(chan WSMessage, 256),
		done:       make(chan struct{}),
	}
}

// === PUBLIC GETTERS (sudah ada, pastikan ada) ===
func (c *Client) GetSessionUID() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.sessionUID
}

func (c *Client) GetUserID() uint {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.userID
}

func (c *Client) GetRole() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.role
}

// === READ/WRITE PUMP (tidak diubah, pastikan ada) ===
func (c *Client) ReadPump() {
	defer c.cleanup()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket read error: %v", err)
			}
			break
		}

		var wsMsg WSMessage
		if err := json.Unmarshal(msg, &wsMsg); err != nil {
			c.sendError("invalid_message", "Failed to parse message: " + err.Error())
			continue
		}

		c.hub.HandleIncomingMessage(c, wsMsg)
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.cleanup()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			var buf bytes.Buffer
			if err := json.NewEncoder(&buf).Encode(message); err != nil {
				log.Printf("Failed to encode message: %v", err)
				continue
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, buf.Bytes()); err != nil {
				log.Printf("Failed to write message: %v", err)
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("Failed to send ping: %v", err)
				return
			}

		case <-c.done:
			return
		}
	}
}

func (c *Client) Send(message WSMessage) {
	select {
	case c.send <- message:
	default:
		log.Printf("Client send buffer full, dropping message for user %d", c.userID)
	}
}

func (c *Client) sendError(code, message string) {
	c.Send(NewWSMessage(WSMsgTypeError, ErrorPayload{
		Code:    code,
		Message: message,
	}))
}

func (c *Client) cleanup() {
	c.once.Do(func() {
		close(c.done)
		c.hub.UnregisterClient(c)
		c.conn.Close()
		close(c.send)
	})
}