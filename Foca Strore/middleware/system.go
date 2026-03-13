package middleware

import (
	"net/http"
	"os"
	"voca-store/response"

	"github.com/gin-gonic/gin"
)

func SystemAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		systemPassword := os.Getenv("SYSTEM_PASSWORD")

		if systemPassword == "" {
			response.ErrorResponse(c, http.StatusInternalServerError, "system password not configured")
			c.Abort()
			return
		}

		password := c.PostForm("password")
		if password == "" {
			response.ErrorResponse(c, http.StatusBadRequest, "password required (form-data)")
			c.Abort()
			return
		}

		if password != systemPassword {
			response.ErrorResponse(c, http.StatusBadRequest, "invalid system password")
			c.Abort()
			return
		}

		c.Next()
	}
}