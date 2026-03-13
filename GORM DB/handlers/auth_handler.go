package handlers

import (
	"main/database"
	"main/helpers"
	"main/models"
	"main/request"
	"main/respons"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Register(c *gin.Context) {
	var req request.RegisterRequest
	c.ShouldBindJSON(&req)

	hashed, err := helpers.HashPassword(req.Password)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			respons.NewJsonResponse("Hash error", nil),
		)
		return
	}

	user := models.User{
		Name: req.Name,
		Email: req.Email,
		Password: hashed,
		RoleID: 2,
	}

	if err := database.DB.Preload("roles").Create(&user).Error; err != nil {
		c.JSON(
			http.StatusBadRequest,
			request.NewJsonResponse("Error", err.Error()),
		)
		return
	}

	if err := database.DB.Preload("Role").First(&user, user.ID).Error; err != nil {
		c.JSON(
			http.StatusInternalServerError,
			respons.NewJsonResponse("Error", err.Error()),
		)
	}

	c.JSON(http.StatusOK, request.NewJsonResponse("Succes", respons.User{
		ID: user.ID,
		Name: user.Name,
		Email: user.Email,
		Role: respons.Role{
			ID: user.Role.ID,
			Name: user.Role.Name,
		},
	}))
}

func Login(c *gin.Context) {
	var req request.LoginRequest
	c.ShouldBindJSON(&req)

	var user models.User
	if err := database.DB.Preload("Role").Where("email = ?", req.Email).First(&user).Error; err != nil{
		c.JSON(
			401,
			respons.NewJsonResponse("Email not found", nil),
		)
		return
	}

	if user.ID == 0 {
		c.JSON(401, respons.NewJsonResponse("User not found", nil))
		return
	}

	if err := helpers.CheckPassword(user.Password, req.Password); err != nil {
		c.JSON(http.StatusForbidden, respons.NewJsonResponse("Wrong Password", nil))
		return
	}

	accessToken, _ := helpers.GenerateAccessToken(user)
	refreshToken, _ := helpers.GenerateRefreshToken(user)

	if err := database.DB.Create(&models.RefreshToken{
		UserID: user.ID,
		Token: refreshToken,
		ExpiredAt: time.Now().Add(
			7 * 24 * time.Hour,
			// 10 * time.Second,
		),
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError,
		respons.NewJsonResponse("Error", err.Error()))
		return
	}

	c.JSON(
		200,
		respons.NewJsonResponse("Login berhasil", respons.TokenResponse{
			AccessToken: accessToken,
			RefreshToken: refreshToken,
			Role: user.Role.Name,
		}),
	)
}

func RefreshToken(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, respons.NewJsonResponse("Request tidak valid", nil))
		return
	}

	token, err := jwt.Parse(body.RefreshToken, func(t *jwt.Token) (interface{}, error) {
		return helpers.REFRESH_SECRET, nil
	})

	var stored models.RefreshToken
	now := time.Now()
	if err := database.DB.Where("token = ? AND expired_at > ?", body.RefreshToken, now).First(&stored).Error; err != nil {
		c.JSON(401, respons.NewJsonResponse("Refresh token sudah logout / tidak valid", nil))
		return
	}

	if err != nil || !token.Valid {
		c.JSON(401, respons.NewJsonResponse("Refresh token tidak valid", nil))
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(401, respons.NewJsonResponse("Invalid token claims", nil))
		return
	}

	var user models.User
	if err := database.DB.Preload("Role").First(&user, claims["user_id"]).Error; err != nil {
		c.JSON(404, respons.NewJsonResponse("User tidak ditemukan", nil))
		return
	}

	newAccess, err := helpers.GenerateAccessToken(user)
	if err != nil {
		c.JSON(500, respons.NewJsonResponse("Gagal membuat access token", nil))
		return
	}

	c.JSON(200, respons.NewJsonResponse(
		"Access token berhasil diperbarui",
		respons.RefreshTokenResponse{
			AccessToken: newAccess,
			Role: user.Role.Name,
		},
	))
}

func Logout(c *gin.Context) {
	refreshToken := c.GetHeader("X-Refresh-Token")

	if refreshToken == "" {
		c.JSON(
			400,
			respons.NewJsonResponse("Refresh token not found", nil),
		)
		return
	}

	result := database.DB.
		Where("token = ?", refreshToken).
		Delete(&models.RefreshToken{})

	if result.Error != nil {
		c.JSON(
			500,
			respons.NewJsonResponse("Failed logout", nil),
		)
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(
			401,
			respons.NewJsonResponse("Refresh token not valid", nil),
		)
		return
	}

	c.JSON(
		http.StatusOK,
		respons.NewJsonResponse("Logout Success", nil),
	)
}

func LogoutAll(c *gin.Context) {
	rt := c.GetHeader("x-Refresh-Token")

	if rt == "" {
		c.JSON(
			http.StatusBadRequest,
			respons.NewJsonResponse("Refresh token not found", nil),
		)
		return
	}

	var current models.RefreshToken
	if err := database.DB.Where("token = ?", rt).First(&current).Error; err != nil {
		c.JSON(
			http.StatusUnauthorized,
			respons.NewJsonResponse("Refresh token not valid", nil),
		)
		return
	}

	result := database.DB.Where("user_id = ? AND token <> ?", current.UserID, rt).Delete(&models.RefreshToken{})

	if result.Error != nil {
		c.JSON(
			http.StatusInternalServerError,
			respons.NewJsonResponse("Failed logout", nil),
		)
		return
	}

	c.JSON(
		http.StatusOK,
		respons.NewJsonResponse("Logout all succes", gin.H{
			"deleted_sessions": result.RowsAffected,
		}),
	)
}