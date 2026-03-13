package models

type User struct {
	ID uint `gorm:"primaryKey"`

	Name string
	Email string `gorm:"unique"`
	Password string `json:"-"`
	
	RoleID uint
	Role Role
	
	Barangs []Barang `gorm:"many2many:user_barangs; constraint:OnDelete:Cascade"`
}