package models

import (
	"time"

	"gorm.io/gorm"
)

type UserRole string

const (
	RoleAdmin UserRole = "admin"
	RoleUser  UserRole = "user"
	RoleMod   UserRole = "mod"
)

type User struct {
	ID        uint          `gorm:"primaryKey" json:"id"`
	Email     string        `gorm:"uniqueIndex;not null" json:"email"`
	Password  string        `gorm:"not null" json:"-"`
	FullName  string        `json:"full_name"`
	Role      UserRole      `gorm:"type:varchar(20);default:'user'" json:"role"`
	IsActive  bool          `gorm:"default:true" json:"is_active"`
	IsBanned  bool          `gorm:"default:false" json:"is_banned"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}