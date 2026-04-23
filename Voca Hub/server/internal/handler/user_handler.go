package handler

import (
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"server/internal/domain/dto"
	"server/internal/helper"
	"server/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Me(c *gin.Context) {
	user := helper.MustCurrentUser(c)
	helper.Success(c, http.StatusOK, "current user fetched", dto.BuildUserResponse(*user))
}

func (h *UserHandler) UpdateMe(c *gin.Context) {
	var req dto.UpdateProfileRequest
	if err := c.ShouldBind(&req); err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	profile, err := c.FormFile("profile")
	if err != nil {
		profile = nil
	}

	username, exists := c.GetPostForm("username")
	if exists {
		req.Username = &username
	}

	user := helper.MustCurrentUser(c)
	response, err := h.userService.UpdateProfileResponse(user.ID, req.Username, profile)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		}
		helper.Error(c, status, err.Error())
		return
	}

	helper.Success(c, http.StatusOK, "profile updated", response)
}

func (h *UserHandler) ServeProfile(c *gin.Context) {
	filePath := strings.TrimPrefix(c.Param("filepath"), "/")
	reader, contentType, err := h.userService.OpenProfileAsset(filePath)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	defer reader.Close()

	if contentType == "" {
		contentType = mime.TypeByExtension(strings.ToLower(filepath.Ext(filePath)))
	}
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	c.Header("Content-Type", contentType)
	c.Header("Cache-Control", "public, max-age=3600")

	if _, err := io.Copy(c.Writer, reader); err != nil {
		c.Status(http.StatusInternalServerError)
	}
}
