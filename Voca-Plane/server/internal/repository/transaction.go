package repository

import (
	"context"
	"voca-plane/internal/domain/models"

	"gorm.io/gorm"
)

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(ctx context.Context, tx *gorm.DB, t *models.Transaction) error {
	if err := tx.WithContext(ctx).Create(t).Error; err != nil {
		return err
	}

	return tx.WithContext(ctx).
		Preload("Items").
		Preload("Items.FlightClass").
		Preload("Flight").
		Preload("Flight.Airline").
		Preload("Flight.Origin").
		Preload("Flight.Destination").
		First(t, t.ID).Error
}

func (r *transactionRepository) GetByCode(ctx context.Context, code string) (*models.Transaction, error) {
	var t models.Transaction
	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Items.FlightClass").
		Preload("Flight").
		Preload("Flight.Airline").
		Preload("Flight.Origin").
		Preload("Flight.Destination").
		First(&t, "code = ?", code).Error
	return &t, err
}

func (r *transactionRepository) GetByUserID(ctx context.Context, userID uint, page, limit int) ([]models.Transaction, int64, error) {
	var transactions []models.Transaction
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Transaction{}).
		Where("user_id = ?", userID).
		Preload("Flight").
		Preload("Flight.Airline").
		Preload("Flight.Origin").
		Preload("Flight.Destination").
		Preload("Items").
		Preload("Items.FlightClass")

	query.Count(&total)
	offset := (page - 1) * limit
	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&transactions).Error
	return transactions, total, err
}

func (r *transactionRepository) GetByUserIDAll(ctx context.Context, userID uint) ([]models.Transaction, error) {
	var transactions []models.Transaction

	query := r.db.WithContext(ctx).Model(&models.Transaction{}).
		Where("user_id = ?", userID).
		Preload("Flight").
		Preload("Flight.Airline").
		Preload("Flight.Origin").
		Preload("Flight.Destination").
		Preload("Items").
		Preload("Items.FlightClass")

	err := query.Order("created_at DESC").Find(&transactions).Error
	return transactions, err
}

func (r *transactionRepository) UpdatePaymentStatus(ctx context.Context, tx *gorm.DB, id uint, status string) error {
	return tx.WithContext(ctx).Model(&models.Transaction{}).Where("id = ?", id).Update("payment_status", status).Error
}

func (r *transactionRepository) Delete(ctx context.Context, tx *gorm.DB, code string) error {
	return tx.WithContext(ctx).Where("code = ?", code).Delete(&models.Transaction{}).Error
}

func (r *transactionRepository) CreateTransactionItems(ctx context.Context, tx *gorm.DB, items []models.TransactionItem) error {
	return tx.WithContext(ctx).Create(&items).Error
}

func (r *transactionRepository) UpdatePaymentURL(ctx context.Context, code string, url string) error {
	return r.db.WithContext(ctx).
		Model(&models.Transaction{}).
		Where("code = ?", code).
		Update("payment_url", url).Error
}
