package dto

import "server/internal/domain/models"

type CategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

type CategoryResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type DifficultyResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func BuildCategoryResponses(categories []models.Category) []CategoryResponse {
	result := make([]CategoryResponse, 0, len(categories))
	for _, category := range categories {
		result = append(result, BuildCategoryResponse(category))
	}
	return result
}

func BuildCategoryResponse(category models.Category) CategoryResponse {
	return CategoryResponse{
		ID:   category.ID,
		Name: category.Name,
	}
}

func BuildDifficultyResponses(difficulties []models.Difficulty) []DifficultyResponse {
	result := make([]DifficultyResponse, 0, len(difficulties))
	for _, difficulty := range difficulties {
		result = append(result, BuildDifficultyResponse(difficulty))
	}
	return result
}

func BuildDifficultyResponse(difficulty models.Difficulty) DifficultyResponse {
	return DifficultyResponse{
		ID:   difficulty.ID,
		Name: difficulty.Name,
	}
}
