package middleware

import (
	"net/http"
	"voca-plane/pkg/response"

	"github.com/gin-gonic/gin"
)

func AppPassword(appPassword string) gin.HandlerFunc {
	return func(c *gin.Context) {
		password := c.GetHeader("password-app")
		if password == "" || password != appPassword {
			response.Error(c, http.StatusUnauthorized, "invalid app password")
			c.Abort()
			return
		}

		c.Next()
	}
}
