package repository

import (
	"server/internal/domain/models"
	domainrepo "server/internal/domain/repository"

	"gorm.io/gorm"
)

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) domainrepo.MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(message *models.Message) error {
	return r.db.Create(message).Error
}

func (r *messageRepository) ListConversation(userID uint, otherUserID uint) ([]models.Message, error) {
	var messages []models.Message
	err := r.db.
		Where(
			"(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
			userID, otherUserID, otherUserID, userID,
		).
		Order("created_at asc").
		Find(&messages).Error
	return messages, err
}
