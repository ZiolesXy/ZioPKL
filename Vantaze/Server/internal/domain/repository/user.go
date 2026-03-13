package repository

import (
	"context"
	"vantaze/internal/domain/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, tx *gorm.DB, user *models.User) error
	GetByID(ctx context.Context, tx *gorm.DB, id uint) (*models.User,  error)
	GetByEmail(ctx context.Context, tx *gorm.DB, email string) (*models.User, error)
	Update(ctx context.Context, tx *gorm.DB, user *models.User) error
	Delete(ctx context.Context, tx *gorm.DB, id uint) error
	List(ctx context.Context, tx *gorm.DB, limit, offset int, search string) ([]models.User, int64, error)
	Ban(ctx context.Context, tx *gorm.DB, id uint) error
	UnBan(ctx context.Context, tx *gorm.DB, id uint) error
}