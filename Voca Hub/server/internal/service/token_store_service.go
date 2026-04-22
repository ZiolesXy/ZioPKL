package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"server/internal/helper"
)

type TokenStoreService struct {
	redis *redis.Client
}

func NewTokenStoreService(redisClient *redis.Client) *TokenStoreService {
	return &TokenStoreService{redis: redisClient}
}

func (s *TokenStoreService) StoreRefreshToken(userID uint, refreshToken string, expiresAt time.Time) error {
	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		return errors.New("refresh token already expired")
	}

	return s.redis.Set(context.Background(), s.refreshTokenKey(userID), helper.HashToken(refreshToken), ttl).Err()
}

func (s *TokenStoreService) ValidateRefreshToken(userID uint, refreshToken string) (bool, error) {
	value, err := s.redis.Get(context.Background(), s.refreshTokenKey(userID)).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return value == helper.HashToken(refreshToken), nil
}

func (s *TokenStoreService) DeleteRefreshToken(userID uint) error {
	return s.redis.Del(context.Background(), s.refreshTokenKey(userID)).Err()
}

func (s *TokenStoreService) BlacklistAccessToken(accessToken string, expiresAt time.Time) error {
	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		return nil
	}

	return s.redis.Set(context.Background(), s.accessBlacklistKey(accessToken), "1", ttl).Err()
}

func (s *TokenStoreService) IsAccessTokenBlacklisted(accessToken string) (bool, error) {
	value, err := s.redis.Exists(context.Background(), s.accessBlacklistKey(accessToken)).Result()
	if err != nil {
		return false, err
	}

	return value > 0, nil
}

func (s *TokenStoreService) refreshTokenKey(userID uint) string {
	return fmt.Sprintf("auth:refresh:%d", userID)
}

func (s *TokenStoreService) accessBlacklistKey(accessToken string) string {
	return "auth:blacklist:" + helper.HashToken(accessToken)
}
