package app

import (
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetUpServer() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	allowed := []string{"http://localhost:3000"}
	if env := os.Getenv("ALLOW_ORIGINS"); env != "" {
		allowed = strings.Split(env, ",")
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins: allowed,
		AllowMethods: []string{"GET", "POST", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders: []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge: 12 * time.Hour,
	}))

	return r
}