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

type AuthMiddleware struct {
	tokenManager *helper.TokenManager
	tokenStore   *service.TokenStoreService
	userService  *service.UserService
}

func NewAuthMiddleware(tokenManager *helper.TokenManager, tokenStore *service.TokenStoreService, userService *service.UserService) *AuthMiddleware {
	return &AuthMiddleware{
		tokenManager: tokenManager,
		tokenStore:   tokenStore,
		userService:  userService,
	}
}

func (m *AuthMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := helper.ExtractBearerToken(c.GetHeader("Authorization"))
		if err != nil && isWebSocketRequest(c) {
			queryToken := strings.TrimSpace(c.Query("token"))
			if queryToken != "" {
				token = queryToken
				err = nil
			}
		}
		if err != nil {
			helper.Error(c, http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		claims, err := m.tokenManager.VerifyAccessToken(token)
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

		blacklisted, err := m.tokenStore.IsAccessTokenBlacklisted(token)
		if err != nil {
			helper.Error(c, http.StatusInternalServerError, err.Error())
			c.Abort()
			return
		}
		if blacklisted {
			helper.Error(c, http.StatusUnauthorized, "token has been revoked")
			c.Abort()
			return
		}

		user, err := m.userService.GetByID(claims.UserID)
		if err != nil {
			helper.Error(c, http.StatusInternalServerError, err.Error())
			c.Abort()
			return
		}
		if user == nil {
			helper.Error(c, http.StatusUnauthorized, "user not found")
			c.Abort()
			return
		}

		if isDebugMode() {
			log.Printf("authenticated user id=%d role=%s email=%s", user.ID, user.Role, user.Email)
		}

		c.Set(helper.ContextClaimsKey, claims)
		c.Set(helper.ContextUserKey, user)
		c.Set(helper.ContextTokenKey, token)
		c.Next()
	}
}

func isWebSocketRequest(c *gin.Context) bool {
	upgrade := strings.TrimSpace(c.GetHeader("Upgrade"))
	connection := strings.TrimSpace(c.GetHeader("Connection"))

	return strings.EqualFold(upgrade, "websocket") || strings.Contains(strings.ToLower(connection), "upgrade")
}

func isDebugMode() bool {
	mode := strings.TrimSpace(os.Getenv("GIN_MODE"))
	return mode == "" || strings.EqualFold(mode, gin.DebugMode)
}
