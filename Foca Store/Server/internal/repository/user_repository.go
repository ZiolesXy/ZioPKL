package repository

import (
	"voca-store/internal/domain/models"
	"voca-store/internal/domain/repository"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{db}
}

func (r *userRepository) FindByEmail(email string) (models.User, error) {
	var user models.User
	err := r.db.Preload("Role").Where("email = ?", email).First(&user).Error
	return user, err
}

func (r *userRepository) FindByID(id uint) (models.User, error) {
	var user models.User
	err := r.db.Preload("Role").Preload("Addresses").First(&user, id).Error
	return user, err
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) Update(user *models.User, updates map[string]interface{}) error {
	return r.db.Model(user).Updates(updates).Error
}

func (r *userRepository) FindRoleByName(name string) (models.Role, error) {
	var role models.Role
	err := r.db.Where("name = ?", name).First(&role).Error
	return role, err
}

func (r *userRepository) CreateCart(cart *models.Cart) error {
	return r.db.Create(cart).Error
}
