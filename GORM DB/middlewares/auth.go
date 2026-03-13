package middlewares

import (
	"main/helpers"
	"main/respons"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401,
				respons.NewJsonResponse("No token provided", nil))
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return helpers.ACCESS_SECRET, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(401,
				respons.NewJsonResponse("Invalid token", nil))
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(401,
				respons.NewJsonResponse("Invalid token claims", nil))
			return
		}

		if claims["type"] != "access" {
			c.AbortWithStatusJSON(401,
				respons.NewJsonResponse("Invalid token type", nil))
			return
		}

		userID := uint(claims["user_id"].(float64))
		role := claims["role"].(string)

		c.Set("user_id", userID)
		c.Set("role", role)

		c.Next()
	}
}
