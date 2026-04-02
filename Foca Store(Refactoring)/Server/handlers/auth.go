package handlers

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"time"
	"voca-store/database"
	"voca-store/helper"
	"voca-store/models"
	"voca-store/request"
	"voca-store/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{DB: db}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	var existingUser models.User
	if err := h.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		response.ErrorResponse(c, http.StatusConflict, "Email already registered")
		return
	}

	hashedPassword, err := helper.HashPassword(req.Password)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	var userRole models.Role
	if err := h.DB.Where("name = ?", "User").First(&userRole).Error; err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "User role not found")
		return
	}

	user := models.User{
		Name:            req.Name,
		Email:           req.Email,
		Password:        hashedPassword,
		TelephoneNumber: req.TelephoneNumber,
		RoleID:          userRole.ID,
	}

	if err := h.DB.Create(&user).Error; err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to create user")
		return
	}

	if err := h.DB.Preload("Role").First(&user, user.ID).Error; err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to load user data")
		return
	}

	cart := models.Cart{UserID: user.ID}
	if err := h.DB.Create(&cart).Error; err != nil {
		response.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"failed to create cart",
		)
		return
	}

	userResp := response.BuildUserResponse(user.ID, user.Name, user.Email, user.Role.Name, user.TelephoneNumber)
	response.SuccessResponse(c, "user registered succesfully", userResp)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "invallid request body")
		return
	}

	var user models.User
	if err := h.DB.Preload("Role").Where("email = ?", req.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.ErrorResponse(c, http.StatusBadRequest, "invalid email or password")
		} else {
			response.ErrorResponse(c, http.StatusInternalServerError, "database error")
		}
		return
	}

	if err := helper.VerifyPassword(user.Password, req.Password); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "invalid email or password")
		return
	}

	accessToken, err := helper.GenerateAccessToken(user.ID, user.Role.Name)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "failed to get access token")
		return
	}

	refreshToken, err := helper.GenerateRefreshToken(user.ID, user.Role.Name)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "failed to get refresh token")
		return
	}

	err = database.RDB.Set(
		database.Ctx,
		"refresh:"+fmt.Sprint(user.ID),
		refreshToken,
		7*24*time.Hour,
	).Err()
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "failed store refresh token")
		return
	}

	userResp := response.BuildUserResponse(user.ID, user.Name, user.Email, user.Role.Name, user.TelephoneNumber)
	authResp := response.BuildAuthResponse(userResp, accessToken, refreshToken)

	response.SuccessResponse(c, "login succesfull", authResp)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req request.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	claims, err := helper.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		response.ErrorResponse(c, http.StatusUnauthorized, "invalid or expired refresh token")
		return
	}

	var user models.User
	if err := h.DB.Preload("Role").First(&user, claims.UserID).Error; err != nil {
		response.ErrorResponse(c, http.StatusNotFound, "user not found")
		return
	}

	key := fmt.Sprintf("refresh:%d", claims.UserID)
	storedToken, err := database.RDB.Get(
		database.Ctx,
		key,
	).Result()

	if err != nil {
		response.ErrorResponse(c, http.StatusUnauthorized, "refresh token not found")
		return
	}

	if storedToken != req.RefreshToken {
		response.ErrorResponse(c, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	accessToken, err := helper.GenerateAccessToken(user.ID, user.Role.Name)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "failed to get access token")
		return
	}

	tokenResp := response.BuildToken(accessToken)
	response.SuccessResponse(c, "token refreshed succesfull", tokenResp)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	UserID := userIDRaw.(uint)

	authHeader := c.GetHeader("Authorization")
	tokenString := authHeader[7:]

	claims, err := helper.ValidateAccessToken(tokenString)
	if err != nil {
		response.ErrorResponse(c, http.StatusUnauthorized, "invalid token")
		return
	}

	exp := claims.ExpiresAt.Time
	ttl := time.Until(exp)

	//blacklist
	err = database.RDB.Set(
		database.Ctx,
		"blacklist:"+tokenString,
		"revoked",
		ttl,
	).Err()

	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "failed blacklist token")
		return
	}

	//delete ref token
	database.RDB.Del(
		database.Ctx,
		"refresh:"+fmt.Sprint(UserID),
	)

	response.SuccessResponse(c, "logout success", nil)
}

func ChangePassword(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDRaw, exist := c.Get("user_id")
		if !exist {
			response.ErrorResponse(c, http.StatusUnauthorized, "missing authorization header")
			return
		}

		userID := userIDRaw.(uint)
		var req request.ChangePasswordRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.ErrorResponse(c, http.StatusBadRequest, "invalid body request")
			return
		}

		var user models.User
		if err := db.First(&user, userID).Error; err != nil {
			response.ErrorResponse(c, http.StatusNotFound, "user not found")
			return
		}

		//check old pw
		if err := helper.VerifyPassword(
			user.Password,
			req.OldPassword,
		); err != nil {
			response.ErrorResponse(c, http.StatusBadRequest, "old password incorect")
			return
		}

		//new pw & conf pw same
		if req.NewPassword != req.ConfirmPassword {
			response.ErrorResponse(c, http.StatusBadRequest, "password confirmation does not match")
			return
		}

		// new pw must not same
		if err := helper.VerifyPassword(
			user.Password,
			req.NewPassword,
		); err == nil {
			response.ErrorResponse(c, http.StatusBadRequest, "new password must be different")
			return
		}

		//hashing new password
		hash, err := helper.HashPassword(
			req.NewPassword,
		)
		if err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "failed hash password")
			return
		}

		// update in database
		if err := db.Model(&user).Update("password", hash).Error; err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "failed update password")
			return
		}

		response.SuccessResponse(c, "password change", nil)
	}
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req request.ForgotPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "invalid body request")
		return
	}

	email := req.Email
	limitKey := "otp:limit:" + email
	otpKey := "otp:" + email

	count, err := database.RDB.Get(database.Ctx, limitKey).Int()
	if err != nil && err.Error() != "redis: nil" {
		response.ErrorResponse(c, http.StatusInternalServerError, "redis error")
		return
	}

	if count >= 5 {
		response.ErrorResponse(c, http.StatusTooManyRequests, "too many otp request, try again later")
		return
	}

	var user models.User
	if err := h.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		database.RDB.Incr(database.Ctx, limitKey)
		database.RDB.Expire(database.Ctx, limitKey, 5*time.Minute)

		response.SuccessResponse(c, "if email exists, otp sent", nil)
		return
	}

	n, _ := rand.Int(rand.Reader, big.NewInt(1000000))
	otp := fmt.Sprintf("%06d", n.Int64())

	err = database.RDB.Set(
		database.Ctx,
		otpKey,
		otp,
		5*time.Minute,
	).Err()

	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "failed generate otp")
		return
	}

	newCount, _ := database.RDB.Incr(database.Ctx, limitKey).Result()
	if newCount == 1 {
		database.RDB.Expire(database.Ctx, limitKey, 5*time.Minute)
	}

	go func() {
		if err := helper.SendOTPEmail(email, otp); err != nil {
			fmt.Println("failed send otp email:", err)
		}
	}()

	response.SuccessResponse(c, "if email exists, otp sent", nil)
}

func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req request.VerifyOTPRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "invalid body request")
		return
	}

	key := "otp:" + req.Email

	storedOTP, err := database.RDB.Get(database.Ctx, key).Result()
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "invalid or expired otp")
		return
	}

	if storedOTP != req.OTP {
		response.ErrorResponse(c, http.StatusBadRequest, "invalid otp")
		return
	}

	hashed, err := helper.HashPassword(req.NewPassword)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "failed hashed password")
		return
	}

	err = h.DB.Model(&models.User{}).
		Where("email = ?", req.Email).
		Update("password", hashed).Error
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "failed update password")
		return
	}

	database.RDB.Del(database.Ctx, key)
	response.SuccessResponse(c, "password reset success", nil)
}
