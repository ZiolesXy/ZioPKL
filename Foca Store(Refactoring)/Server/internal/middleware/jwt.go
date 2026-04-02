package middleware

import (
	"fmt"
	"net/http"
	"voca-store/internal/database"
	"voca-store/internal/helper"
	"voca-store/internal/domain/dto/response"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func JWTAuth(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.ErrorResponse(c, http.StatusUnauthorized, "Missing authorization header")
			c.Abort()
			return
		}

		tokenString := authHeader
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		}

		claims, err := helper.ValidateAccessToken(tokenString)
		if err != nil {
			response.ErrorResponse(c, http.StatusUnauthorized, "Invalid or expired token")
			c.Abort()
			return
		}

		// Ambil data dari blacklist
		_, err = database.RDB.Get(database.Ctx, "blacklist:"+tokenString).Result()

		if err == nil {
			// Token ditemukan di blacklist -> Blokir!
			response.ErrorResponse(c, http.StatusUnauthorized, "token revoked")
			c.Abort()
			return
		} else if err != redis.Nil {
			// Error di sini BUKAN karena data kosong, tapi karena hal lain (misal: Redis MATI)

			fmt.Println("Gagal koneksi ke Redis:", err)

			// Jika Anda ingin SANGAT AMAN (Redis mati = tidak boleh akses):
			response.ErrorResponse(c, http.StatusInternalServerError, "auth service error")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)

		c.Next()
	}
}
