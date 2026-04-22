package dto

import (
	"time"

	"server/internal/domain/models"
)

type CreatePostRequest struct {
	Content string `json:"content" binding:"required"`
}

type UpdatePostRequest struct {
	Content string `json:"content" binding:"required"`
}

type PostResponse struct {
	ID        uint         `json:"id"`
	UserID    uint         `json:"user_id"`
	Content   string       `json:"content"`
	CreatedAt string       `json:"created_at"`
	User      UserResponse `json:"user"`
}

func BuildPostResponses(posts []models.Post) []PostResponse {
	result := make([]PostResponse, 0, len(posts))
	for _, post := range posts {
		result = append(result, BuildPostResponse(post))
	}
	return result
}

func BuildPostResponse(post models.Post) PostResponse {
	return PostResponse{
		ID:        post.ID,
		UserID:    post.UserID,
		Content:   post.Content,
		CreatedAt: post.CreatedAt.Format(time.RFC3339),
		User:      BuildUserResponse(post.User),
	}
}
