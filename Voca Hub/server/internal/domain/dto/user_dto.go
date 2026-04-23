package dto

import (
	"net/url"
	"strings"
	"time"

	"server/internal/domain/models"
)

type UserResponse struct {
	ID         uint    `json:"id"`
	Email      string  `json:"email"`
	Username   *string `json:"username"`
	ProfileURL *string `json:"profile_url"`
	Role       string  `json:"role"`
	CreatedAt  string  `json:"created_at,omitempty"`
}

type UpdateProfileRequest struct {
	Username *string `form:"username"`
}

func BuildUserResponses(users []models.User) []UserResponse {
	result := make([]UserResponse, 0, len(users))
	for _, user := range users {
		result = append(result, BuildUserResponse(user))
	}
	return result
}

func BuildUserResponse(user models.User) UserResponse {
	profileURL := buildProfileURL(user.ProfileURL)

	response := UserResponse{
		ID:         user.ID,
		Email:      user.Email,
		Username:   user.Username,
		ProfileURL: profileURL,
		Role:       user.Role,
	}

	if !user.CreatedAt.IsZero() {
		response.CreatedAt = user.CreatedAt.Format(time.RFC3339)
	}

	return response
}

func buildProfileURL(profilePath *string) *string {
	if profilePath == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*profilePath)
	if trimmed == "" {
		return nil
	}
	if strings.Contains(trimmed, "://") {
		parsed, err := url.Parse(trimmed)
		if err == nil {
			trimmed = parsed.Path
		}
	}

	url := "/users/profile/" + strings.TrimLeft(trimmed, "/")
	return &url
}
