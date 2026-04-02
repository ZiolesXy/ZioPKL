package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID                   uint   `gorm:"primaryKey"`
	Name                 string `gorm:"not null"`
	Email                string `gorm:"unique;not null"`
	Password             string `gorm:"not null"`
	TelephoneNumber      string `gorm:"type:varchar(20)"`
	ProfileImageURL      string `gorm:"type:text"`
	ProfileImagePublicID string `gorm:"type:text"`
	RoleID               uint   `gorm:"not null"`
	Role                 Role   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            gorm.DeletedAt
	Checkouts            []Checkout `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Addresses            []Address  `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
