package models

type Difficulty struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Name  string `gorm:"size:100;not null;uniqueIndex" json:"name"`
	Games []Game `json:"-"`
}
