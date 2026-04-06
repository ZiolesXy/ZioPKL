package repository

import (
	"voca-store/internal/domain/models"
)

type ProductRepository interface {
	Create(product *models.Product) error
	FindAll() ([]models.Product, error)
	FindBySlug(slug string) (models.Product, error)
	FindByID(id uint) (models.Product, error)
	Update(product *models.Product, updates map[string]interface{}) error
	Delete(product *models.Product) error
	DeleteAll() error
	ExistsBySlug(slug string) (bool, error)
	FindWithImages() ([]models.Product, error)
	ClearAllImages() error
}
