package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"server/internal/domain/dto"
	"server/internal/helper"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (h *UserHandler) Me(c *gin.Context) {
	user := helper.MustCurrentUser(c)
	helper.Success(c, http.StatusOK, "current user fetched", dto.BuildUserResponse(*user))
}
