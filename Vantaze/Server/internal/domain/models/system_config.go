package models

import "time"

type SystemConfig struct {
	ID        uint     `gorm:"primaryKey" json:"id"`
	Key       string   `gorm:"uniqueIndex;not null" json:"key"`
	Value     string   `gorm:"type:text" json:"value"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (SystemConfig) TableName() string {
	return "system_configs"
}