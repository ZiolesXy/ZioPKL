package repository

import "server/internal/domain/models"

type UserRepository interface {
	FindByID(id uint) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Create(user *models.User) error
	Save(user *models.User) error
	List() ([]models.User, error)
	CountAll() (int64, error)
	CountByRole(role string) (int64, error)
}
