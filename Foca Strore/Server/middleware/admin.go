package middleware

import (
	"net/http"
	"voca-store/response"

	"github.com/gin-gonic/gin"
)

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exist := c.Get("role")
		if !exist || role != "Admin" {
			response.ErrorResponse(c, http.StatusForbidden, "access denied. admin previlage require")
			c.Abort()
			return
		}
		c.Next()
	}
}
