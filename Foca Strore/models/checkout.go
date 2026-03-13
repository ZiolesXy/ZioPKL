package models

import (
	"time"

	"gorm.io/gorm"
)

type Checkout struct {
	ID     uint   `gorm:"primaryKey"`
	UID    string `gorm:"uniqueIndex;not null"`
	UserID uint   `gorm:"not null"`
	User   *User  `gorm:"foreignKey:UserID"`

	CouponID *uint
	Coupon   *Coupon `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	AddressID *uint
	Address   *Address `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	Subtotal       float64 `gorm:"not null;default:0"`
	DiscountAmount float64 `gorm:"not null;default:0"`

	TotalPrice float64 `gorm:"not null"`
	Status     string  `gorm:"not null;type:varchar(20);default:'pending'"`

	MidtransOrderID string `gorm:"uniqueIndex"`
	SnapToken       string
	PaymentURL      string
	PaymentStatus   string `gorm:"type:varchar(20);default:'pending'"`

	Items       []CheckoutItem `gorm:"foreignKey:CheckoutID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	WhatsappURL string         `gorm:"type:text"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}