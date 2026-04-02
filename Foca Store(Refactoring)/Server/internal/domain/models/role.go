package models

// import "gorm.io/gorm"

type Role struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique;not null"`
}

// func (r *Role) BeforeCreate(tx *gorm.DB) error {
// 	r.Name = r.Name
// 	return nil
// }
