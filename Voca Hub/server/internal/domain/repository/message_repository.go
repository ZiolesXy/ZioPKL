package repository

import "server/internal/domain/models"

type MessageRepository interface {
	Create(message *models.Message) error
	ListConversation(userID uint, otherUserID uint) ([]models.Message, error)
}
