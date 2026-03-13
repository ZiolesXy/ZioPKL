package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)
		log.Printf("[%d] %s %s %v", c.Writer.Status(), c.Request.Method, c.Request.URL.Path, latency)
	}
}