package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"
	"voca-store/internal/domain/repository"

	"github.com/redis/go-redis/v9"
)

type authRepository struct {
	rdb *redis.Client
	ctx context.Context
}

func NewAuthRepository(rdb *redis.Client) repository.AuthRepository {
	return &authRepository{rdb, context.Background()}
}

func (r *authRepository) SetRefreshToken(userID uint, token string, ttl time.Duration) error {
	key := fmt.Sprintf("refresh:%d", userID)
	return r.rdb.Set(r.ctx, key, token, ttl).Err()
}

func (r *authRepository) GetRefreshToken(userID uint) (string, error) {
	key := fmt.Sprintf("refresh:%d", userID)
	return r.rdb.Get(r.ctx, key).Result()
}

func (r *authRepository) DeleteRefreshToken(userID uint) error {
	key := fmt.Sprintf("refresh:%d", userID)
	return r.rdb.Del(r.ctx, key).Err()
}

func (r *authRepository) BlacklistToken(token string, ttl time.Duration) error {
	key := "blacklist:" + token
	return r.rdb.Set(r.ctx, key, "revoked", ttl).Err()
}

func (r *authRepository) IsTokenBlacklisted(token string) (bool, error) {
	key := "blacklist:" + token
	val, err := r.rdb.Exists(r.ctx, key).Result()
	return val > 0, err
}

func (r *authRepository) SetOTP(email string, otp string, ttl time.Duration) error {
	key := "otp:" + email
	return r.rdb.Set(r.ctx, key, otp, ttl).Err()
}

func (r *authRepository) GetOTP(email string) (string, error) {
	key := "otp:" + email
	return r.rdb.Get(r.ctx, key).Result()
}

func (r *authRepository) DeleteOTP(email string) error {
	key := "otp:" + email
	return r.rdb.Del(r.ctx, key).Err()
}

func (r *authRepository) IncrementOTPLimit(email string, ttl time.Duration) (int64, error) {
	key := "otp:limit:" + email
	n, err := r.rdb.Incr(r.ctx, key).Result()
	if n == 1 {
		r.rdb.Expire(r.ctx, key, ttl)
	}
	return n, err
}

func (r *authRepository) GetOTPLimit(email string) (int, error) {
	key := "otp:limit:" + email
	val, err := r.rdb.Get(r.ctx, key).Result()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(val)
}
