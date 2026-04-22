package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"server/internal/domain/dto"
	"server/internal/helper"
	"server/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response, err := h.authService.Register(req)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "email already registered" {
			status = http.StatusConflict
		}
		helper.Error(c, status, err.Error())
		return
	}

	helper.Success(c, http.StatusCreated, "register success", response)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response, err := h.authService.Login(req)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "invalid email or password" {
			status = http.StatusUnauthorized
		}
		helper.Error(c, status, err.Error())
		return
	}

	helper.Success(c, http.StatusOK, "login success", response)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response, err := h.authService.Refresh(req.RefreshToken)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "invalid refresh token" {
			status = http.StatusUnauthorized
		}
		helper.Error(c, status, err.Error())
		return
	}

	helper.Success(c, http.StatusOK, "refresh success", response)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	user := helper.MustCurrentUser(c)
	claims := c.MustGet(helper.ContextClaimsKey).(*helper.AuthClaims)
	accessToken := c.MustGet(helper.ContextTokenKey).(string)

	if claims.ExpiresAt == nil {
		helper.Error(c, http.StatusBadRequest, "missing access token expiry")
		return
	}

	if err := h.authService.Logout(user.ID, accessToken, claims.ExpiresAt.Time); err != nil {
		helper.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.Success(c, http.StatusOK, "logout success", nil)
}
