package models

type Role struct {
	ID uint `gorm:"primaryKey"`
	Name string

	Users []User `gorm:"foreignKey:RoleID; constraint:OnUpdate:CASCADE, OnDelete:SET NULL;" json:"users,omitempty"`
}