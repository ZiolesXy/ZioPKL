package models

import "time"

type CartItem struct {
	ID        uint    `gorm:"primaryKey"`
	CartID    uint    `gorm:"not null"`
	ProductID uint    `gorm:"not null"`
	Product   *Product `gorm:"foreignKey:ProductID"`
	Quantity  int
	CreatedAt time.Time
	UpdatedAt time.Time
}