package repository

import (
	"voca-store/internal/domain/models"
)

type CategoryRepository interface {
	Create(category *models.Category) error
	FindAll() ([]models.Category, map[uint]int64, error)
	FindBySlug(slug string) (models.Category, []models.Product, error)
	FindByID(id uint) (models.Category, error)
	Update(category *models.Category, updates map[string]interface{}) error
	Delete(category *models.Category) error
	ExistsBySlug(slug string) (bool, error)
}
