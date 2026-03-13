package models

type Barang struct {
	ID uint `gorm:"primaryKey"`
	Name string

	Users []User`gorm:"many2many:user_barangs; constraint:OnDelete:Cascade" json:"users,omitempty"`
}