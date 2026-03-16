package repository

import (
	"context"
	"voca-plane/internal/domain/models"
	"voca-plane/pkg/helper"

	"gorm.io/gorm"
)

type promoRepository struct {
	db *gorm.DB
}

func NewPromoRepository(db *gorm.DB) PromoRepository {
	return &promoRepository{db: db}
}

func (r *promoRepository) GetByID(ctx context.Context, id uint) (*models.PromoCode, error) {
	var promo models.PromoCode
	err := r.db.WithContext(ctx).First(&promo, id).Error
	return &promo, err
}

func (r *promoRepository) GetByCode(ctx context.Context, code string) (*models.PromoCode, error) {
	var promo models.PromoCode
	err := r.db.WithContext(ctx).Where("code = ? AND is_active = ?", code, true).First(&promo).Error
	return &promo, err
}

func (r *promoRepository) GetAll(ctx context.Context, page, limit int, sortBy, order string) ([]models.PromoCode, int64, error) {
	var promos []models.PromoCode
	var total int64

	query := r.db.WithContext(ctx).Model(&models.PromoCode{})
	query.Session(&gorm.Session{}).Count(&total)

	// Promos Whitelist
	allowedColumns := map[string]bool{
		"id":        true,
		"code":      true,
		"discount":  true,
		"is_active": true,
	}

	query = helper.ApplySorting(query, sortBy, order, allowedColumns, "id ASC")

	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Find(&promos).Error
	return promos, total, err
}

func (r *promoRepository) Create(ctx context.Context, tx *gorm.DB, promo *models.PromoCode) error {
	return tx.WithContext(ctx).Create(promo).Error
}

func (r *promoRepository) Update(ctx context.Context, tx *gorm.DB, promo *models.PromoCode) error {
	return tx.WithContext(ctx).Save(promo).Error
}

func (r *promoRepository) Delete(ctx context.Context, tx *gorm.DB, id uint) error {
	return tx.WithContext(ctx).Delete(&models.PromoCode{}, id).Error
}