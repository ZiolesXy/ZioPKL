package repository

import (
	"context"
	"voca-plane/internal/domain/models"

	"gorm.io/gorm"
)

type systemRepository struct {
	db *gorm.DB
}

func NewSystemRepository(db *gorm.DB) SystemRepository {
	return &systemRepository{db: db}
}

func (r *systemRepository) ResetDatabase(ctx context.Context) error {
	// Drop all tables
	err := r.db.Migrator().DropTable(
		&models.User{},
		&models.Airline{},
		&models.Airport{},
		&models.Flight{},
		&models.FlightClass{},
		&models.Seat{},
		&models.FlightSeat{},
		&models.Transaction{},
		&models.TransactionItem{},
		&models.PromoCode{},
	)
	if err != nil {
		return err
	}

	// Re-migrate tables
	return r.db.AutoMigrate(
		&models.User{},
		&models.Airline{},
		&models.Airport{},
		&models.Flight{},
		&models.FlightClass{},
		&models.Seat{},
		&models.FlightSeat{},
		&models.Transaction{},
		&models.TransactionItem{},
		&models.PromoCode{},
	)
}
