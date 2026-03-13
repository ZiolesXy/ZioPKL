package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func CORS(allowedOrigins string) gin.HandlerFunc {
    return func(c *gin.Context) {
        origin := c.Request.Header.Get("Origin")
        
        // Set Origin secara dinamis atau wildcard
        if allowedOrigins == "*" {
            c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        } else if origin != "" && strings.Contains(allowedOrigins, origin) {
            c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
        }

        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        
        // PENTING: Gabungkan semua header yang diizinkan dalam SATU baris
        allowedHeaders := "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, ngrok-skip-browser-warning"
        c.Writer.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
        
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

        // Handle Preflight Request
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(http.StatusNoContent)
            return
        }

        c.Next()
    }
}
