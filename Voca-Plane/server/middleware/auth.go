package middleware

import (
	"errors"
	"net/http"
	"strings"
	"voca-plane/internal/repository"
	"voca-plane/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth(secret string, userRepo repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "missing authorization header")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, http.StatusUnauthorized, "invalid authorization format")
			c.Abort()
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			response.Error(c, http.StatusUnauthorized, "invalid token")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			response.Error(c, http.StatusUnauthorized, "invalid token claims")
			c.Abort()
			return
		}

		if claims["type"] != "access" {
			response.Error(c, http.StatusUnauthorized, "invalid token type")
			c.Abort()
			return
		}

		userIDFloat, ok := claims["id"].(float64)
		if !ok {
			response.Error(c, http.StatusUnauthorized, "invalid user id")
			c.Abort()
			return
		}

		userID := uint(userIDFloat)

		// 🔥 cek database user
		user, err := userRepo.FindByID(c.Request.Context(), userID)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "user not found")
			c.Abort()
			return
		}

		if user.IsBanned {
			response.Error(c, http.StatusForbidden, "account banned")
			c.Abort()
			return
		}

		email, _ := claims["email"].(string)
		role, _ := claims["role"].(string)

		c.Set("userID", userID)
		c.Set("userRole", role)
		c.Set("userEmail", email)

		c.Next()
	}
}