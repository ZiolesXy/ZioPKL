package dto

import (
	"time"

	"server/internal/domain/models"
)

type UserResponse struct {
	ID         uint    `json:"id"`
	ClerkID    string  `json:"clerk_id"`
	Email      string  `json:"email"`
	Username   *string `json:"username"`
	ProfileURL *string `json:"profile_url"`
	Role       string  `json:"role"`
	CreatedAt  string  `json:"created_at,omitempty"`
}

func BuildUserResponses(users []models.User) []UserResponse {
	result := make([]UserResponse, 0, len(users))
	for _, user := range users {
		result = append(result, BuildUserResponse(user))
	}
	return result
}

func BuildUserResponse(user models.User) UserResponse {
	response := UserResponse{
		ID:         user.ID,
		ClerkID:    user.ClerkID,
		Email:      user.Email,
		Username:   user.Username,
		ProfileURL: user.ProfileURL,
		Role:       user.Role,
	}

	if !user.CreatedAt.IsZero() {
		response.CreatedAt = user.CreatedAt.Format(time.RFC3339)
	}

	return response
}
