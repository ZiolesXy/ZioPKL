package models

import "time"

type Coupon struct {
	ID    uint    `gorm:"primaryKey"`
	Code  string  `gorm:"uniqueIndex;not null"`
	Type  string  `gorm:"type:varchar(20);not null"`
	Value float64 `gorm:"not null"`
	Quota int     `gorm:"not null"`

	UsedCount int `gorm:"not null"`
	MinimumPurchase float64 `gorm:"type:decimal(15,2);default:0"`

	IsActive  *bool `gorm:"default:true"`
	ExpiresAt *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}
