package service

import (
	"context"
	"errors"
	"fmt"
	"time"
	"voca-plane/internal/domain/dto/request"
	"voca-plane/internal/domain/models"
	"voca-plane/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo           repository.UserRepository
	jwtSecret          string
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string, accessTokenExpiry, refreshTokenExpiry time.Duration) *AuthService {
	return &AuthService{
		userRepo:           userRepo,
		jwtSecret:          jwtSecret,
		accessTokenExpiry:  accessTokenExpiry,
		refreshTokenExpiry: refreshTokenExpiry,
	}
}

func(s *AuthService) Register(ctx context.Context, req request.RegisterRequest) error {
	_, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err == nil {
		return errors.New("email already registerd")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &models.User{
		Name: req.Name,
		Email: req.Email,
		Password: string(hashedPassword),
		Role: models.RoleUser,
	}

	return s.userRepo.Create(ctx, user)
}

func (s *AuthService) Login(ctx context.Context, req request.LoginRequest) (*models.TokenPair, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if user.IsBanned != false {
		return nil, errors.New("you are banned")
	}

	if user.DeletedAt.Valid {
		return nil, errors.New("this account was deleted")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	accessToken, err := s.generateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateRefreshToken(user.ID, user.Role)

	return &models.TokenPair{
		AccessToken: accessToken,
		RefreshToken: refreshToken,
		ExpiresIn: int(s.accessTokenExpiry.Seconds()),
		TokenType: "Bearer",
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*models.TokenPair, error) {
	claims, err := s.validateRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	userID := uint(claims["id"].(float64))
	// userRole := claims["role"].(string)

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	accessToken, err := s.generateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	return &models.TokenPair{
		AccessToken: accessToken,
		RefreshToken: refreshToken,
		ExpiresIn: int(s.accessTokenExpiry.Seconds()),
		TokenType: "Bearer",
	}, nil
}

func(s *AuthService) generateAccessToken(userID uint, email, role string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = userID
	claims["email"] = email
	claims["role"] = role
	claims["type"] = "access"
	claims["exp"] = time.Now().Add(s.accessTokenExpiry).Unix()
	claims["iat"] = time.Now().Unix()

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	return tokenString, err
}

func (s *AuthService) generateRefreshToken(userID uint, role string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = userID
	claims["role"] = role
	claims["type"] = "refresh"
	claims["jti"] = uuid.New().String()
	claims["exp"] = time.Now().Add(s.refreshTokenExpiry).Unix()

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	return tokenString, err
}

func (s *AuthService) validateRefreshToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	if claims["type"] != "refresh" {
		return nil, errors.New("invalid token type")
	}

	return claims, nil
}

func buildRefreshTokenKey(userID uint, token string) string {
	return fmt.Sprintf("refresh_token:%d:%s", userID, token)
}