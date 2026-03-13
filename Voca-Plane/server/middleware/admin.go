package middleware

import (
	"net/http"
	"voca-plane/pkg/response"

	"github.com/gin-gonic/gin"
)

func RequireAdmin() gin.HandlerFunc{
	return func(c *gin.Context) {
		userRole, exists := c.Get("userRole")
		if !exists {
			response.Error(c, http.StatusUnauthorized, "role not found")
			c.Abort()
			return
		}

		role := userRole.(string)
		if role != "admin" && role != "super_admin" {
			response.Error(c, http.StatusUnauthorized, "admin access required")
			c.Abort()
			return
		}

		c.Next()
	}
}

func RequireSuperAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("userRole")
		if !exists {
			response.Error(c, http.StatusForbidden, "role not found")
			c.Abort()
			return
		}

		role := userRole.(string)
		if role != "super_admin" {
			response.Error(c, http.StatusForbidden, "super admin access required")
			c.Abort()
			return
		}

		c.Next()
	}
}