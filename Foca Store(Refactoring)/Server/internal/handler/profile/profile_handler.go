package profile

import (
	"net/http"
	"voca-store/internal/domain/dto/request"
	"voca-store/internal/domain/dto/response"
	"voca-store/internal/service"

	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	authService service.AuthService
}

func NewProfileHandler(authService service.AuthService) *ProfileHandler {
	return &ProfileHandler{authService}
}

func (h *ProfileHandler) GetProfile(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uint)

	profileRes, err := h.authService.GetProfile(userID)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "failed to get profile")
		return
	}

	response.SuccessResponse(c, "profile retrieved", profileRes)
}

func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uint)

	var req request.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	profileRes, err := h.authService.UpdateProfile(userID, req)
	if err != nil {
		response.ErrorResponse(c, http.StatusConflict, err.Error())
		return
	}

	response.SuccessResponse(c, "profile updated", profileRes)
}

func (h *ProfileHandler) ChangePassword(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uint)

	var req request.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.authService.ChangePassword(userID, req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessResponse(c, "password changed", nil)
}
