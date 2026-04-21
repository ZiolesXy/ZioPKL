package repository

import (
	"errors"

	"gorm.io/gorm"

	"server/internal/domain/models"
	domainrepo "server/internal/domain/repository"
)

type gameRepository struct {
	db *gorm.DB
}

func NewGameRepository(db *gorm.DB) domainrepo.GameRepository {
	return &gameRepository{db: db}
}

func (r *gameRepository) Create(game *models.Game) error {
	return r.db.Create(game).Error
}

func (r *gameRepository) FindByID(id uint) (*models.Game, error) {
	var game models.Game
	err := r.db.Preload("Developer").First(&game, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &game, nil
}

func (r *gameRepository) Update(game *models.Game) error {
	return r.db.Save(game).Error
}

func (r *gameRepository) ListAll() ([]models.Game, error) {
	var games []models.Game
	err := r.db.Preload("Developer").Order("created_at desc").Find(&games).Error
	return games, err
}

func (r *gameRepository) ListApproved() ([]models.Game, error) {
	var games []models.Game
	err := r.db.Preload("Developer").Where("status = ?", "approved").Order("created_at desc").Find(&games).Error
	return games, err
}

func (r *gameRepository) ListByDeveloperID(developerID uint) ([]models.Game, error) {
	var games []models.Game
	err := r.db.Preload("Developer").Where("developer_id = ?", developerID).Order("created_at desc").Find(&games).Error
	return games, err
}

func (r *gameRepository) CountAll() (int64, error) {
	var total int64
	err := r.db.Model(&models.Game{}).Count(&total).Error
	return total, err
}
