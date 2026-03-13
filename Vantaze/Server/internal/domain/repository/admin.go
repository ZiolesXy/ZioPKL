package repository

import (
	"context"
	"vantaze/internal/domain/models"

	"gorm.io/gorm"
)

type AdminRepository interface {
	Create(ctx context.Context, tx *gorm.DB, admin *models.Admin) error
	GetByID(ctx context.Context, tx *gorm.DB, id uint) (*models.Admin, error)
	GetByEmail(ctx context.Context, tx *gorm.DB, email string) (*models.Admin, error)
	Update(ctx context.Context, tx *gorm.DB, admin *models.Admin) error
	Delete(ctx context.Context, tx *gorm.DB, id uint) error
	List(ctx context.Context, tx *gorm.DB, limit, offset int) ([]models.Admin, int64, error)
}