package response

import (
	"time"
	"voca-store/internal/domain/models"
)

type ChatSessionResponse struct {
	UID           string     `json:"uid"`
	UserID        uint       `json:"user_id"`
	UserName      string     `json:"user_name"`
	UserAvatar    string     `json:"user_avatar,omitempty"`
	AdminID       *uint      `json:"admin_id,omitempty"`
	AdminName     *string    `json:"admin_name,omitempty"`
	AdminAvatar   *string    `json:"admin_avatar,omitempty"`
	Status        string     `json:"status"`
	LastMessage   *string    `json:"last_message,omitempty"`
	LastMessageAt *time.Time `json:"last_message_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type ChatMessageResponse struct {
	UID         string     `json:"uid"`
	SessionUID  string     `json:"session_uid"`
	SenderID    uint       `json:"sender_id"`
	SenderName  string     `json:"sender_name"`
	SenderRole  string     `json:"sender_role"`
	Content     string     `json:"content"`
	MessageType string     `json:"message_type"`
	MediaURL    *string    `json:"media_url,omitempty"`
	ReadAt      *time.Time `json:"read_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

type WebSocketMessage struct {
	Type      string      `json:"type"`
	Payload   interface{} `json:"payload"`
	Timestamp time.Time   `json:"timestamp"`
}

func ToChatSessionResponse(session models.ChatSession) ChatSessionResponse {
	resp := ChatSessionResponse{
		UID: session.UID,
		UserID: session.UserID,
		UserName: session.User.Name,
		UserAvatar: session.User.ProfileImageURL,
		AdminID: session.AdminID,
		Status: string(session.Status),
		LastMessage: session.LastMessage,
		LastMessageAt: session.LastMessageAt,
		CreatedAt: session.CreatedAt,
		UpdatedAt: session.UpdatedAt,
	}
	if session.AdminID != nil {
		resp.AdminName = &session.Admin.Name
		resp.AdminAvatar = &session.Admin.ProfileImageURL
	}
	return resp
}

func ToChatMessageResponse(msg models.ChatMessage, senderRole string) ChatMessageResponse {
	return ChatMessageResponse{
		UID: msg.UID,
		SessionUID: msg.Session.UID,
		SenderID: msg.SenderID,
		SenderName: msg.Sender.Name,
		SenderRole: msg.Sender.Role.Name,
		Content: msg.Content,
		MessageType: string(msg.MessageType),
		MediaURL: msg.MediaURL,
		ReadAt: msg.ReadAt,
		CreatedAt: msg.CreatedAt,
	}
}