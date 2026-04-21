package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"server/internal/helper"
)

type RoleMiddleware struct{}

func NewRoleMiddleware() *RoleMiddleware {
	return &RoleMiddleware{}
}

func (m *RoleMiddleware) Require(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := helper.MustCurrentUser(c)
		for _, role := range roles {
			if user.Role == role {
				c.Next()
				return
			}
		}
		helper.Error(c, http.StatusForbidden, "forbidden")
		c.Abort()
	}
}
