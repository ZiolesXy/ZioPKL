package models

import (
	"time"

	"gorm.io/gorm"
)

type MessageType string

const (
	MesssageText  MessageType = "text"
	MesssageImage MessageType = "image"
	MesssageFile  MessageType = "file"
)

type ChatMessage struct {
	ID          uint        `gorm:"primaryKey"`
	UID         string      `gorm:"type:uuid;default:gen_random_uuid();uniqueIndex;not null"`
	SessionID   uint        `gorm:"not null;index"`
	Session     ChatSession `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	SenderID    uint        `gorm:"not null;index"`
	Sender      User        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Content     string      `gorm:"type:text;not null"`
	MessageType MessageType `gorm:"type:varchar(20);default:'text'"`
	MediaURL    *string     `gorm:"type:text"`
	ReadAt      *time.Time   `gorm:"index"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}