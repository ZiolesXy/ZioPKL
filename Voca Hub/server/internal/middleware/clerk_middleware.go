package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	"server/internal/helper"
	"server/internal/service"
)

type ClerkMiddleware struct {
	verifier    *helper.ClerkVerifier
	userService *service.UserService
}

func NewClerkMiddleware(verifier *helper.ClerkVerifier, userService *service.UserService) *ClerkMiddleware {
	return &ClerkMiddleware{
		verifier:    verifier,
		userService: userService,
	}
}

func (m *ClerkMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := helper.ExtractBearerToken(c.GetHeader("Authorization"))
		if err != nil {
			helper.Error(c, http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		claims, err := m.verifier.VerifyToken(token)
		if err != nil {
			log.Printf("token verification failed: %v", err)

			message := "invalid token"
			if isDebugMode() {
				message = err.Error()
			}

			helper.Error(c, http.StatusUnauthorized, message)
			c.Abort()
			return
		}

		user, err := m.userService.SyncUser(*claims)
		if err != nil {
			helper.Error(c, http.StatusInternalServerError, err.Error())
			c.Abort()
			return
		}

		if isDebugMode() {
			log.Printf(
				"authenticated user id=%d clerk_id=%s role=%s email=%s",
				user.ID,
				user.ClerkID,
				user.Role,
				user.Email,
			)
		}

		c.Set(helper.ContextClaimsKey, claims)
		c.Set(helper.ContextUserKey, user)
		c.Next()
	}
}

func isDebugMode() bool {
	mode := strings.TrimSpace(os.Getenv("GIN_MODE"))
	return mode == "" || strings.EqualFold(mode, gin.DebugMode)
}
