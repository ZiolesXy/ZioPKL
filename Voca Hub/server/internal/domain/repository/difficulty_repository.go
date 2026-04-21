package repository

import "server/internal/domain/models"

type DifficultyRepository interface {
	FindByID(id uint) (*models.Difficulty, error)
	ListAll() ([]models.Difficulty, error)
}
