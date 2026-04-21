package models

type Category struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Name  string `gorm:"size:100;not null;uniqueIndex" json:"name"`
	Games []Game `gorm:"many2many:game_categories;" json:"-"`
}
