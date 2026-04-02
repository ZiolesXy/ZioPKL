package repository

import (
	"voca-store/internal/domain/models"
)

type CouponRepository interface {
	Create(coupon *models.Coupon) error
	FindAll() ([]models.Coupon, error)
	FindByID(id uint) (models.Coupon, error)
	Update(coupon *models.Coupon, updates map[string]interface{}) error
	Delete(coupon *models.Coupon) error
	
	// User Coupon
	FindUserCoupon(userID uint, couponID uint) (models.UserCoupon, error)
	FindUserCouponByID(id uint, userID uint) (models.UserCoupon, error)
	CreateUserCoupon(userCoupon *models.UserCoupon) error
	FindMyCoupons(userID uint) ([]models.UserCoupon, error)
	DeleteUserCoupon(userCoupon *models.UserCoupon) error
	
	// Transaction
	WithTransaction(fn func(repo CouponRepository) error) error
	IncrementUsedCount(couponID uint, amount int) error
}
