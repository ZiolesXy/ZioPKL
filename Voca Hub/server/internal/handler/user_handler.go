package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"server/internal/helper"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (h *UserHandler) Me(c *gin.Context) {
	user := helper.MustCurrentUser(c)
	helper.Success(c, http.StatusOK, "current user fetched", gin.H{
		"id":          user.ID,
		"clerk_id":    user.ClerkID,
		"email":       user.Email,
		"username":    user.Username,
		"profile_url": user.ProfileURL,
		"role":        user.Role,
	})
}
