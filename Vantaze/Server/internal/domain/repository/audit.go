package repository

import (
	"context"
	"vantaze/internal/domain/models"

	"gorm.io/gorm"
)

type AuditRepository interface {
	Log(ctx context.Context, tx *gorm.DB, log *models.AuditLog) error
}