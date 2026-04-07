package middleware

import (
	"net/http"
	"voca-store/internal/database"
	"voca-store/internal/domain/dto/response"
	"voca-store/internal/helper"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func WebSocketAuth(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Query("token")
		if tokenString == "" {
			response.ErrorResponse(c, http.StatusUnauthorized, "missing token")
			c.Abort()
			return
		}

		claims, err := helper.ValidateAccessToken(tokenString)
		if err != nil {
			response.ErrorResponse(c, http.StatusUnauthorized, "invalid or expired token")
			c.Abort()
			return
		}

		// Check Redis blacklist
		_, err = database.RDB.Get(database.Ctx, "blacklist:"+tokenString).Result()
		if err == nil {
			response.ErrorResponse(c, http.StatusUnauthorized, "token revoked")
			c.Abort()
			return
		} else if err != redis.Nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "auth service error")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}