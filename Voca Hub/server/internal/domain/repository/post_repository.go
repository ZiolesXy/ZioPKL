package repository

import "server/internal/domain/models"

type PostRepository interface {
	Create(post *models.Post) error
	FindByID(id uint) (*models.Post, error)
	ListAll() ([]models.Post, error)
	ListByUserID(userID uint) ([]models.Post, error)
	Update(post *models.Post) error
	Delete(id uint) error
}
