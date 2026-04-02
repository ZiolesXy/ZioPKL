package repository

import (
	"voca-store/internal/domain/models"
)

type AddressRepository interface {
	Create(address *models.Address) error
	FindByID(id uint, userID uint) (models.Address, error)
	FindByUID(uid string, userID uint) (models.Address, error)
	FindAll(userID uint) ([]models.Address, error)
	Update(address *models.Address, updates map[string]interface{}) error
	Delete(address *models.Address) error
	
	UnsetPrimary(userID uint) error
	WithTransaction(fn func(repo AddressRepository) error) error
}
