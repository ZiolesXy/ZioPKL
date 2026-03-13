package models

import "time"

type Cart struct {
	ID        uint       `gorm:"primaryKey"`
	UserID    uint       `gorm:"unique;not null"`
	User      *User       `gorm:"foreignKey:UserID"`
	Items     []CartItem `gorm:"foreignKey:CartID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedAt time.Time
	UpdatedAt time.Time
}