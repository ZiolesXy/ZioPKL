package service

import (
	"errors"
	"strings"
	"time"

	"server/internal/domain/dto"
	"server/internal/domain/models"
	domainrepo "server/internal/domain/repository"
	"server/internal/helper"
)

type AuthService struct {
	userRepo     domainrepo.UserRepository
	tokenManager *helper.TokenManager
	tokenStore   *TokenStoreService
}

func NewAuthService(userRepo domainrepo.UserRepository, tokenManager *helper.TokenManager, tokenStore *TokenStoreService) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		tokenManager: tokenManager,
		tokenStore:   tokenStore,
	}
}

func (s *AuthService) Register(req dto.RegisterRequest) (*dto.AuthResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	existingUser, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	passwordHash, err := helper.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:        email,
		PasswordHash: passwordHash,
		Username:     normalizeOptionalString(req.Username),
		ProfileURL:   normalizeOptionalString(req.ProfileURL),
		Role:         "USER",
	}
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return s.issueTokens(user)
}

func (s *AuthService) Login(req dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(strings.ToLower(strings.TrimSpace(req.Email)))
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	if err := helper.ComparePassword(user.PasswordHash, req.Password); err != nil {
		return nil, errors.New("invalid email or password")
	}

	return s.issueTokens(user)
}

func (s *AuthService) Refresh(refreshToken string) (*dto.AuthResponse, error) {
	refreshToken = strings.TrimSpace(refreshToken)
	claims, err := s.tokenManager.VerifyRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	valid, err := s.tokenStore.ValidateRefreshToken(user.ID, refreshToken)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, errors.New("invalid refresh token")
	}

	return s.issueAccessTokenFromRefresh(user, refreshToken, claims.ExpiresAt.Time)
}

func (s *AuthService) issueTokens(user *models.User) (*dto.AuthResponse, error) {
	tokens, err := s.tokenManager.GenerateTokenPair(user)
	if err != nil {
		return nil, err
	}

	refreshClaims, err := s.tokenManager.VerifyRefreshToken(tokens.RefreshToken)
	if err != nil {
		return nil, err
	}
	if refreshClaims.ExpiresAt == nil {
		return nil, errors.New("missing refresh token expiry")
	}

	if err := s.tokenStore.StoreRefreshToken(user.ID, tokens.RefreshToken, refreshClaims.ExpiresAt.Time); err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		User:   dto.BuildUserResponse(*user),
		Tokens: *tokens,
	}, nil
}

func (s *AuthService) issueAccessTokenFromRefresh(user *models.User, refreshToken string, refreshExpiresAt time.Time) (*dto.AuthResponse, error) {
	tokens, err := s.tokenManager.GenerateTokenPair(user)
	if err != nil {
		return nil, err
	}

	tokens.RefreshToken = refreshToken
	tokens.RefreshTokenExpiresAt = refreshExpiresAt.Format(time.RFC3339)

	return &dto.AuthResponse{
		User:   dto.BuildUserResponse(*user),
		Tokens: *tokens,
	}, nil
}

func (s *AuthService) Logout(userID uint, accessToken string, accessExpiresAt time.Time) error {
	if err := s.tokenStore.DeleteRefreshToken(userID); err != nil {
		return err
	}

	return s.tokenStore.BlacklistAccessToken(accessToken, accessExpiresAt)
}

func normalizeOptionalString(value *string) *string {
	if value == nil {
		return nil
	}

	normalized := strings.TrimSpace(*value)
	if normalized == "" {
		return nil
	}

	return &normalized
}
