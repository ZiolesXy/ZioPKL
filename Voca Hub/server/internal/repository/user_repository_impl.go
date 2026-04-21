package repository

import (
	"errors"

	"gorm.io/gorm"

	"server/internal/domain/models"
	domainrepo "server/internal/domain/repository"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domainrepo.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByClerkID(clerkID string) (*models.User, error) {
	var user models.User
	err := r.db.Where("clerk_id = ?", clerkID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) List() ([]models.User, error) {
	var users []models.User
	err := r.db.Order("created_at desc").Find(&users).Error
	return users, err
}

func (r *userRepository) CountAll() (int64, error) {
	var total int64
	err := r.db.Model(&models.User{}).Count(&total).Error
	return total, err
}

func (r *userRepository) CountByRole(role string) (int64, error) {
	var total int64
	err := r.db.Model(&models.User{}).Where("role = ?", role).Count(&total).Error
	return total, err
}
