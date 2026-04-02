package auth

import (
	"net/http"
	// "strings"
	"voca-store/internal/domain/dto/request"
	"voca-store/internal/domain/dto/response"
	"voca-store/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	userRes, err := h.authService.Register(req)
	if err != nil {
		response.ErrorResponse(c, http.StatusConflict, err.Error())
		return
	}

	response.SuccessResponse(c, "user registered successfully", userRes)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	authRes, err := h.authService.Login(req)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessResponse(c, "login successful", authRes)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req request.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	tokenRes, err := h.authService.RefreshToken(req)
	if err != nil {
		response.ErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	response.SuccessResponse(c, "token refreshed successful", tokenRes)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uint)

	authHeader := c.GetHeader("Authorization")
	if len(authHeader) < 7 {
		response.ErrorResponse(c, http.StatusUnauthorized, "invalid authorization header")
		return
	}
	token := authHeader[7:]

	h.authService.Logout(userID, token)
	response.SuccessResponse(c, "logout success", nil)
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req request.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "invalid body request")
		return
	}

	h.authService.ForgotPassword(req)
	response.SuccessResponse(c, "if email exists, otp sent", nil)
}

func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req request.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "invalid body request")
		return
	}

	if err := h.authService.VerifyOTP(req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessResponse(c, "password reset success", nil)
}
