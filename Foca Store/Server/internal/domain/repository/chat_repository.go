package repository

import (
	"context"
	"voca-store/internal/domain/models"
)

type ChatRepository interface {
	// Session
	CreateSession(ctx context.Context, session *models.ChatSession) error
	GetSessionByUID(ctx context.Context, uid string) (*models.ChatSession, error)
	GetSessionByID(ctx context.Context, id uint) (*models.ChatSession, error)
	GetPendingSessions(ctx context.Context) ([]models.ChatSession, error)
	GetUserActiveSession(ctx context.Context, userID uint) (*models.ChatSession, error)
	UpdateSessionStatus(ctx context.Context, id uint, status models.ChatSessionStatus, adminID *uint) error
	UpdateLastMessage(ctx context.Context, sessionID uint, message string) error
	CloseSession(ctx context.Context, id uint) error

	// Message
	CreateMessage(ctx context.Context, message *models.ChatMessage) error
	GetMessagesBySession(ctx context.Context, sessionID uint, limit, offset int) ([]models.ChatMessage, error)
	GetUnreadMessages(ctx context.Context, sessionID uint, userID uint) ([]models.ChatMessage, error)
	MarkMessagesAsRead(ctx context.Context, sessionID uint, userID uint) error
}