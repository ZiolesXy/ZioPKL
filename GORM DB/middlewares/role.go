package middlewares

import (
	"main/respons"
	"github.com/gin-gonic/gin"
)

func OnlyAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleValue, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(
				401,
				respons.NewJsonResponse("Unauthorized", nil),
			)
			return
		}

		role, ok := roleValue.(string)
		if !ok || role != "Admin" {
			c.AbortWithStatusJSON(
				403,
				respons.NewJsonResponse("Forbidden", nil),
			)
			return
		}

		c.Next()
	}
}
