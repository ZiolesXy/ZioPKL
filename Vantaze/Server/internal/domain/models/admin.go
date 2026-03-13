package models

import (
	"time"

	"gorm.io/gorm"
)

type Admin struct {
	ID        uint          `gorm:"primaryKey" json:"id"`
	Email     string        `gorm:"uniqueIndex;not null" json:"email"`
	Password  string        `gorm:"not null" json:"-"`
	FullName  string        `json:"full_name"`
	Role      UserRole      `gorm:"type:varchar(20);default:'admin'" json:"role"`
	IsActive  bool          `gorm:"default:true" json:"is_active"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

func (Admin) TableName() string {
	return "admins"
}