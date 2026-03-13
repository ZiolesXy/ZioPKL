package models

import "time"

type RefreshToken struct {
	ID uint `gorm:"primaryKey"`
	UserID uint
	Token string `gorm:"unique"`
	ExpiredAt time.Time
	CreatedAt time.Time
}