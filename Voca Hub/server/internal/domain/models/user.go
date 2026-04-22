package models

import "time"

type User struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ClerkID    string    `gorm:"uniqueIndex;size:255;not null" json:"clerk_id"`
	Email      string    `gorm:"size:255" json:"email"`
	Username   *string   `gorm:"size:255" json:"username"`
	ProfileURL *string   `gorm:"size:512" json:"profile_url"`
	Role       string    `gorm:"size:50;not null;default:USER" json:"role"`
	CreatedAt  time.Time `json:"created_at"`
}
