package models

import "time"

type Game struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	Title         string     `gorm:"size:255;not null" json:"title"`
	Description   string     `gorm:"type:text" json:"description"`
	FileURL       string     `gorm:"size:500;not null" json:"file_url"`
	ThumbnailPath string     `gorm:"size:500" json:"-"`
	DeveloperID   uint       `gorm:"not null;index" json:"developer_id"`
	Status        string     `gorm:"size:50;not null;default:pending" json:"status"`
	CreatedAt     time.Time  `json:"created_at"`
	Developer     User       `gorm:"foreignKey:DeveloperID" json:"developer"`
	Categories    []Category `gorm:"many2many:game_categories;" json:"categories"`
}
