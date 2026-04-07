package websocket

import "time"

type WSMessageType string

const (
	WSMsgTypeHandshake    WSMessageType = "handshake"
	WSMsgTypeHandshakeAck WSMessageType = "handshake_ack"
	WSMsgTypeMessage      WSMessageType = "message"
	WSMsgTypeMessageAck   WSMessageType = "message_ack"
	WSMsgTypeTyping       WSMessageType = "typing"
	WSMsgTypeNotification WSMessageType = "notification"
	WSMsgTypeError        WSMessageType = "error"
	WSMsgTypeSessionClose WSMessageType = "session_close"
)

type WSMessage struct {
	Type      WSMessageType `json:"type"`
	Payload   interface{}   `json:"payload"`
	Timestamp time.Time      `json:"timestamp"`
}

type HandshakePayload struct {
	SessionUID string `json:"session_uid"`
	UserID     uint   `json:"user_id"`
	Role       string `json:"role"`
}

type MessagePayload struct {
	SessionUID  string    `json:"session_uid"`
	SenderID    uint      `json:"sender_id"`
	SenderName  string    `json:"sender_name"`
	SenderRole  string    `json:"sender_role"`
	Content     string    `json:"content"`
	MessageType string    `json:"message_type"`
	MediaURL    *string   `json:"media_url,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type NotificationPayload struct {
	Title string      `json:"title"`
	Body  string      `json:"body"`
	Data  interface{} `json:"data,omitempty"`
}

type ErrorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewWSMessage(msgType WSMessageType, payload interface{}) WSMessage {
	return WSMessage{
		Type: msgType,
		Payload: payload,
		Timestamp: time.Now(),
	}
}