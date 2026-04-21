package repository

import "server/internal/domain/models"

type GameRepository interface {
	Create(game *models.Game) error
	FindByID(id uint) (*models.Game, error)
	Update(game *models.Game) error
	ListAll() ([]models.Game, error)
	ListApproved() ([]models.Game, error)
	ListByDeveloperID(developerID uint) ([]models.Game, error)
	CountAll() (int64, error)
}
