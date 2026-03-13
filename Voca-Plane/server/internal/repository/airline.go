package repository

import (
	"context"
	"voca-plane/internal/domain/models"
	"voca-plane/pkg/helper"

	"gorm.io/gorm"
)

type airlineRepository struct {
	db *gorm.DB
}

func NewAirlineRepository(db *gorm.DB) AirlineRepository {
	return &airlineRepository{db: db}
}

func (r *airlineRepository) GetAll(ctx context.Context, page, limit int, sortBy, order string) ([]models.Airline, int64, error) {
	var airlines []models.Airline
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Airline{})
	query.Count(&total)

	// Airlines Whitelist
	allowedColumns := map[string]bool{
		"id":   true,
		"name": true,
		"code": true,
	}

	query = helper.ApplySorting(query, sortBy, order, allowedColumns, "id ASC")

	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Find(&airlines).Error
	return airlines, total, err
}

func (r *airlineRepository) GetByID(ctx context.Context, id uint) (*models.Airline, error) {
	var airline models.Airline
	err := r.db.WithContext(ctx).First(&airline, id).Error
	return &airline, err
}

func (r *airlineRepository) Create(ctx context.Context, tx *gorm.DB, airline *models.Airline) error {
	if err := tx.WithContext(ctx).Create(airline).Error; err != nil {
		return err
	}
	return tx.WithContext(ctx).First(airline, airline.ID).Error
}

func (r *airlineRepository) Update(ctx context.Context, tx *gorm.DB, airline *models.Airline) error {
	if err := tx.WithContext(ctx).Save(airline).Error; err != nil {
		return err
	}
	return tx.WithContext(ctx).First(airline, airline.ID).Error
}

func (r *airlineRepository) Delete(ctx context.Context, tx *gorm.DB, id uint) error {
	return tx.WithContext(ctx).Delete(&models.Airline{}, id).Error
}