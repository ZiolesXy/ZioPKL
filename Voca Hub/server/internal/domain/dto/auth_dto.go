package dto

import "server/internal/helper"

type RegisterRequest struct {
	Email      string  `json:"email" binding:"required,email"`
	Password   string  `json:"password" binding:"required,min=8"`
	Username   *string `json:"username"`
	ProfileURL *string `json:"profile_url"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type AuthResponse struct {
	User   UserResponse     `json:"user"`
	Tokens helper.TokenPair `json:"tokens"`
}
