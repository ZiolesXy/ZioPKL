package models

import "time"

type Product struct {
	ID            uint   `gorm:"primaryKey"`
	Name          string `gorm:"not null"`
	Slug          string `gorm:"uniqueIndex;not null"`
	Description   string
	ImageURL      string
	ImagePublicID string
	Price         float64 `gorm:"not null"`
	Stock         int
	CategoryID    uint
	Category      *Category `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	CartItems     []CartItem `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
