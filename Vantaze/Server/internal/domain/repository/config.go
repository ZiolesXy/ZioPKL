package repository

import (
	"context"
	"vantaze/internal/domain/models"

	"gorm.io/gorm"
)

type ConfigRepository interface {
	GetByKey(ctx context.Context, tx *gorm.DB, key string) (*models.SystemConfig, error)
	Set(ctx context.Context, tx *gorm.DB, cfg *models.SystemConfig) error
}