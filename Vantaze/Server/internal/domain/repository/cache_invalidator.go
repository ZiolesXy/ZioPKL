package repository

import (
	"context"
	"time"
)

type CacheInvalidator interface {
	DeleteByPattern(ctx context.Context, pattern string) error
	DeleteKeys(ctx context.Context, keys ...string) error
	BlacklistToken(ctx context.Context, token string, expiration time.Duration) error
	IsTokenBlacklisted(ctx context.Context, token string) (bool, error)
}