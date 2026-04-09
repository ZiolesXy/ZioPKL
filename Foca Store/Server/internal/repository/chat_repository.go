package repository

import (
	"context"
	"time"
	"voca-store/internal/domain/models"
	"voca-store/internal/domain/repository"

	"gorm.io/gorm"
)

type chatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) repository.ChatRepository {
	return &chatRepository{db: db}
}

func (r *chatRepository) CreateSession(ctx context.Context, session *models.ChatSession) error {
	return r.db.WithContext(ctx).Create(session).Error
}

func (r *chatRepository) GetSessionByUID(ctx context.Context, uid string) (*models.ChatSession, error) {
	var session models.ChatSession
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Admin").
		Where("uid = ? AND deleted_at IS NULL", uid).
		First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *chatRepository) GetSessionByName(ctx context.Context, name string) ([]models.ChatSession, error) {
	var sessions []models.ChatSession

	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Admin").
		Joins("JOIN users u ON u.id = chat_sessions.user_id").
		Where("LOWER(u.name) LIKE LOWER(?)", "%"+name+"%").
		Where("chat_sessions.deleted_at IS NULL").
		Order("chat_sessions.last_message_at DESC").
		Find(&sessions).Error

	return sessions, err
}

func (r *chatRepository) GetSessionByID(ctx context.Context, id uint) (*models.ChatSession, error) {
	var session models.ChatSession
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Admin").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *chatRepository) GetPendingSessions(ctx context.Context) ([]models.ChatSession, error) {
	var sessions []models.ChatSession
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("(status = ? OR status = ?) AND admin_id IS NULL AND deleted_at IS NULL", 
			models.ChatSessionPending, models.ChatSessionActive).
		Order("created_at ASC").
		Find(&sessions).Error
	return sessions, err
}

func (r *chatRepository) GetAllSessions(ctx context.Context) ([]models.ChatSession, error) {
	var sessions []models.ChatSession
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Admin").
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Find(&sessions).Error
	return sessions, err
}

func (r *chatRepository) GetUserActiveSession(ctx context.Context, userID uint) (*models.ChatSession, error) {
	var session models.ChatSession
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Admin").
		Where("user_id = ? AND status = ? AND deleted_at IS NULL", userID, models.ChatSessionActive).
		First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *chatRepository) UpdateSessionStatus(ctx context.Context, id uint, status models.ChatSessionStatus, adminID *uint) error {
	updates := map[string]interface{}{"status": status}
	if adminID != nil {
		updates["admin_id"] = adminID
	}
	return r.db.WithContext(ctx).
		Model(&models.ChatSession{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func (r *chatRepository) UpdateLastMessage(ctx context.Context, sessionID uint, message string) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&models.ChatSession{}).
		Where("id = ?", sessionID).
		Updates(map[string]interface{}{
			"last_message":    message,
			"last_message_at": &now,
		}).Error
}

func (r *chatRepository) IsAdmin(ctx context.Context, userID uint) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Table("users").
		Joins("JOIN roles ON users.role_id = roles.id").
		Where("users.id = ? AND roles.name = ?", userID, "Admin").
		Count(&count).Error
	
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *chatRepository) CloseSession(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Model(&models.ChatSession{}).
		Where("id = ?", id).
		Update("status", models.ChatSessionClosed).Error
}

// === MESSAGE METHODS ===

func (r *chatRepository) CreateMessage(ctx context.Context, message *models.ChatMessage) error {
	return r.db.WithContext(ctx).Create(message).Error
}

func (r *chatRepository) GetMessagesBySession(ctx context.Context, sessionID uint, limit, offset int) ([]models.ChatMessage, error) {
	var messages []models.ChatMessage
	err := r.db.WithContext(ctx).
		Preload("Sender").
		Preload("Sender.Role").
		Preload("Session").
		Where("session_id = ? AND deleted_at IS NULL", sessionID).
		Order("created_at ASC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error
	return messages, err
}

func (r *chatRepository) GetUnreadMessages(ctx context.Context, sessionID uint, userID uint) ([]models.ChatMessage, error) {
	var messages []models.ChatMessage
	err := r.db.WithContext(ctx).
		Preload("Sender").
		Where("session_id = ? AND sender_id != ? AND read_at IS NULL AND deleted_at IS NULL", 
			sessionID, userID).
		Order("created_at ASC").
		Find(&messages).Error
	return messages, err
}

func (r *chatRepository) MarkMessagesAsRead(ctx context.Context, sessionID uint, userID uint) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&models.ChatMessage{}).
		Where("session_id = ? AND sender_id != ? AND read_at IS NULL", sessionID, userID).
		Update("read_at", &now).Error
}