package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	ID uint `json:"id"`
	Email string `json:"email"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func NewAccessClaims(id uint, email, role string, expiration time.Duration) *CustomClaims {
	return &CustomClaims{
		ID: id,
		Email: email,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
}

func NewRefreshClaims(id uint, expiration time.Duration) *CustomClaims {
	return &CustomClaims{
		ID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}
}