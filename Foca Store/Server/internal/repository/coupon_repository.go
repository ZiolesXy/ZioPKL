package repository

import (
	"voca-store/internal/domain/models"
	"voca-store/internal/domain/repository"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type couponRepository struct {
	db *gorm.DB
}

func NewCouponRepository(db *gorm.DB) repository.CouponRepository {
	return &couponRepository{db}
}

func (r *couponRepository) Create(coupon *models.Coupon) error {
	return r.db.Create(coupon).Error
}

func (r *couponRepository) FindAll() ([]models.Coupon, error) {
	var coupons []models.Coupon
	err := r.db.Order("id ASC").Find(&coupons).Error
	return coupons, err
}

func (r *couponRepository) FindByID(id uint) (models.Coupon, error) {
	var coupon models.Coupon
	err := r.db.Clauses(clause.Locking{Strength: "UPDATE"}).First(&coupon, id).Error
	return coupon, err
}

func (r *couponRepository) Update(coupon *models.Coupon, updates map[string]interface{}) error {
	return r.db.Model(coupon).Updates(updates).Error
}

func (r *couponRepository) Delete(coupon *models.Coupon) error {
	return r.db.Delete(coupon).Error
}

func (r *couponRepository) FindUserCoupon(userID uint, couponID uint) (models.UserCoupon, error) {
	var userCoupon models.UserCoupon
	err := r.db.Where("user_id = ? AND coupon_id = ?", userID, couponID).First(&userCoupon).Error
	return userCoupon, err
}

func (r *couponRepository) FindUserCouponByID(id uint, userID uint) (models.UserCoupon, error) {
	var userCoupon models.UserCoupon
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&userCoupon).Error
	return userCoupon, err
}

func (r *couponRepository) CreateUserCoupon(userCoupon *models.UserCoupon) error {
	return r.db.Create(userCoupon).Error
}

func (r *couponRepository) FindMyCoupons(userID uint) ([]models.UserCoupon, error) {
	var userCoupons []models.UserCoupon
	err := r.db.Preload("Coupon").Where("user_id = ?", userID).Order("created_at DESC").Find(&userCoupons).Error
	return userCoupons, err
}

func (r *couponRepository) DeleteUserCoupon(userCoupon *models.UserCoupon) error {
	return r.db.Delete(userCoupon).Error
}

func (r *couponRepository) WithTransaction(fn func(repo repository.CouponRepository) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return fn(NewCouponRepository(tx))
	})
}

func (r *couponRepository) IncrementUsedCount(couponID uint, amount int) error {
	return r.db.Model(&models.Coupon{}).Where("id = ?", couponID).Update("used_count", gorm.Expr("used_count + ?", amount)).Error
}
