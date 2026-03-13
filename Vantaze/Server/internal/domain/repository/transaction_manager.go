package repository

import (
	"context"

	"gorm.io/gorm"
)

type TransactionManager interface {
	ExecuteInTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error
}