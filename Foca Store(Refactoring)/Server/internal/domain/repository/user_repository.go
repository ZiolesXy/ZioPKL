package repository

import (
	"voca-store/internal/domain/models"
)

type UserRepository interface {
	FindByEmail(email string) (models.User, error)
	FindByID(id uint) (models.User, error)
	Create(user *models.User) error
	Update(user *models.User, updates map[string]interface{}) error
	FindRoleByName(name string) (models.Role, error)
	CreateCart(cart *models.Cart) error
}
