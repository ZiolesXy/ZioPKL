package models

import (
	"time"

	"gorm.io/gorm"
)

type ChatSessionStatus string

const (
	ChatSessionPending ChatSessionStatus = "pending"
	ChatSessionActive  ChatSessionStatus = "active"
	ChatSessionClosed  ChatSessionStatus = "closed"
)

type ChatSession struct {
	ID           uint              `gorm:"primaryKey"`
	UID          string            `gorm:"type:uuid;default:gen_random_uuid();uniqueIndex;not null"`
	UserID       uint              `gorm:"not null;index"`
	User         User              `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	AdminID      *uint             `gorm:"index"`
	Admin        User              `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Status       ChatSessionStatus `gorm:"type:varchar(20);default:'pending';index"`
	LastMessage  *string           `gorm:"type:text"`
	LastMessageAt *time.Time       `gorm:"index"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	Messages     []ChatMessage  `gorm:"foreignKey:SessionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (ChatSession) TableName() string {
	return "chat_sessions"
}