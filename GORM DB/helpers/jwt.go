package helpers

import (
	"main/models"
	"time"
	"github.com/golang-jwt/jwt/v5"
)

var ACCESS_SECRET = []byte("PENGENDALI_PETIR")
var REFRESH_SECRET = []byte("PKL_DELAMETA")

func GenerateAccessToken(user models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role": user.Role.Name,
		"exp": time.Now().Add(15 * time.Minute).Unix(),
		"type": "access",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(ACCESS_SECRET)
}

func GenerateRefreshToken(user models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
		"type": "refresh",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(REFRESH_SECRET)
}