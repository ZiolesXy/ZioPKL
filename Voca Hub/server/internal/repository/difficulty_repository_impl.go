package repository

import (
	"errors"

	"gorm.io/gorm"

	"server/internal/domain/models"
	domainrepo "server/internal/domain/repository"
)

type difficultyRepository struct {
	db *gorm.DB
}

func NewDifficultyRepository(db *gorm.DB) domainrepo.DifficultyRepository {
	return &difficultyRepository{db: db}
}

func (r *difficultyRepository) FindByID(id uint) (*models.Difficulty, error) {
	var difficulty models.Difficulty
	err := r.db.First(&difficulty, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &difficulty, nil
}

func (r *difficultyRepository) ListAll() ([]models.Difficulty, error) {
	var difficulties []models.Difficulty
	err := r.db.Order("id asc").Find(&difficulties).Error
	return difficulties, err
}
