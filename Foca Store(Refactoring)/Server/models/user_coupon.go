package models

import "time"

type UserCoupon struct {
	ID       uint   `gorm:"primaryKey"`
	UserID   uint   `gorm:"uniqueIndex:idx_user_coupon;index;not null"`
	User     User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CouponID uint   `gorm:"uniqueIndex:idx_user_coupon;index;not null"`
	Coupon   Coupon `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	UsedAt    *time.Time
	CreatedAt time.Time
}
