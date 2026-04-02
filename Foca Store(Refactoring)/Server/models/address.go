package models

import (
	"time"

	"gorm.io/gorm"
)

type Address struct {
	ID     uint   `gorm:"primaryKey"`
	UID    string `gorm:"uniqueIndex;not null"`
	UserID uint   `gorm:"index;not null"`
	User   User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	Label         string `gorm:"not null"`
	RecipientName string `gorm:"not null"`
	Phone         string `gorm:"not null"`
	AddressLine   string `gorm:"not null"`
	City          string `gorm:"not null"`
	Province      string `gorm:"not null"`
	PostalCode    string `gorm:"not null"`
	IsPrimary     bool   `gorm:"default:false"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
