package handler

import (
	"net/http"
	"voca-plane/internal/domain/dto/request"
	"voca-plane/internal/service"
	"voca-plane/pkg/helper"
	"voca-plane/pkg/response"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	user, err := h.service.GetProfile(c.Request.Context(), userID.(uint))
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "profile retrieved", user)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req request.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.service.UpdateProfile(c.Request.Context(), userID.(uint), req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "profile updated", nil)
}

func (h *UserHandler) GetDeviceInfo(c *gin.Context) {

	clientIP := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	serverInfo := helper.GetServerInfo()
	clientInfo := helper.GetClientInfo(clientIP, userAgent)

	response.Success(c, http.StatusOK, "device retrieved successfully", helper.DeviceDetails{
		Server: serverInfo,
		Client: clientInfo,
	})
}