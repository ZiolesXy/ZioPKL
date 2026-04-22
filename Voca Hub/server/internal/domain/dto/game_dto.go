package dto

import (
	"time"

	"server/internal/domain/models"
)

type UploadGameRequest struct {
	Title        string `form:"title" binding:"required"`
	Description  string `form:"description"`
	CategoryIDs  []uint `form:"category_id"`
	DifficultyID uint   `form:"difficulty_id" binding:"required"`
}

type UpdateGameRequest struct {
	Title        string `form:"title"`
	Description  string `form:"description"`
	CategoryIDs  []uint `form:"category_id"`
	DifficultyID *uint  `form:"difficulty_id"`
}

type GameResponse struct {
	ID           uint              `json:"id"`
	Title        string            `json:"title"`
	Description  string            `json:"description"`
	FileURL      string            `json:"file_url"`
	ThumbnailURL string            `json:"thumbnail_url"`
	DeveloperID  uint              `json:"developer_id"`
	Status       string            `json:"status"`
	CreatedAt    string            `json:"upload_at"`
	Developer    models.User       `json:"developer"`
	Difficulty   models.Difficulty `json:"difficulty"`
	Categories   []models.Category `json:"categories"`
}

func BuildGameResponses(games []models.Game, buildThumbnailURL func(string) string) []GameResponse {
	result := make([]GameResponse, 0, len(games))
	for _, game := range games {
		result = append(result, BuildGameResponse(&game, buildThumbnailURL))
	}
	return result
}

func BuildGameResponse(game *models.Game, buildThumbnailURL func(string) string) GameResponse {
	thumbnailURL := ""
	if buildThumbnailURL != nil {
		thumbnailURL = buildThumbnailURL(game.ThumbnailPath)
	}

	return GameResponse{
		ID:           game.ID,
		Title:        game.Title,
		Description:  game.Description,
		FileURL:      game.FileURL,
		ThumbnailURL: thumbnailURL,
		DeveloperID:  game.DeveloperID,
		Status:       game.Status,
		CreatedAt:    game.CreatedAt.Format(time.RFC3339),
		Developer:    game.Developer,
		Difficulty:   game.Difficulty,
		Categories:   game.Categories,
	}
}
