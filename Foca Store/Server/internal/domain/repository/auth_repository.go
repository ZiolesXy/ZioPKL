package repository

import (
	"time"
)

type AuthRepository interface {
	SetRefreshToken(userID uint, token string, ttl time.Duration) error
	GetRefreshToken(userID uint) (string, error)
	DeleteRefreshToken(userID uint) error
	BlacklistToken(token string, ttl time.Duration) error
	IsTokenBlacklisted(token string) (bool, error)
	SetOTP(email string, otp string, ttl time.Duration) error
	GetOTP(email string) (string, error)
	DeleteOTP(email string) error
	IncrementOTPLimit(email string, ttl time.Duration) (int64, error)
	GetOTPLimit(email string) (int, error)
}
