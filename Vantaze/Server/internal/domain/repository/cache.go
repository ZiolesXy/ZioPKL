package repository

import (
	"context"
	"time"
)

type CacheRepository interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string, dest interface{}) error
	Delete(ctx context.Context, key string) error
	Incr(ctx context.Context, key string) (int64, error)
	Expire(ctx context.Context, key string, expiration time.Duration) error
	Exists(ctx context.Context, key string) (bool, error)
}