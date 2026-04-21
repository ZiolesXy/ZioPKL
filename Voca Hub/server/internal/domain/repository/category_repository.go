package repository

import "server/internal/domain/models"

type CategoryRepository interface {
	Create(category *models.Category) error
	FindByID(id uint) (*models.Category, error)
	ListAll() ([]models.Category, error)
	FindByIDs(ids []uint) ([]models.Category, error)
	Update(category *models.Category) error
	Delete(id uint) error
}
