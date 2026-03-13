package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secretKey string
}

func NewJWTManager(secretKey string) *JWTManager {
	return &JWTManager{secretKey: secretKey}
}

func (m *JWTManager) GenerateAccessToken(id uint, email, role string) (string, error) {
	claims := NewAccessClaims(id, email, role, 1*time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(m.secretKey))
}

func (m *JWTManager) GenerateRefreshToken(id uint) (string, error) {
	claims := NewRefreshClaims(id, 24*7*time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(m.secretKey))
}

func (m *JWTManager)  ValidateToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || claims == nil {
		return nil, errors.New("invalid claims type")
	}

	return claims, nil
}

func (m *JWTManager) GetTokenExpiration(tokenString string) (time.Time, error) {
	claims, err := m.ValidateToken(tokenString)
	if err != nil {
		return time.Time{}, err
	}

	if claims.ExpiresAt == nil {
		return time.Time{}, errors.New("np expiration claims")
	}

	return claims.ExpiresAt.Time, nil
}