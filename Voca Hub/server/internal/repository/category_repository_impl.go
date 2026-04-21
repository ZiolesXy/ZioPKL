package repository

import (
	"errors"

	"gorm.io/gorm"

	"server/internal/domain/models"
	domainrepo "server/internal/domain/repository"
)

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) domainrepo.CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(category *models.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) FindByID(id uint) (*models.Category, error) {
	var category models.Category
	err := r.db.First(&category, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) ListAll() ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Order("name asc").Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) FindByIDs(ids []uint) ([]models.Category, error) {
	var categories []models.Category
	if len(ids) == 0 {
		return categories, nil
	}
	err := r.db.Where("id IN ?", ids).Order("name asc").Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) Update(category *models.Category) error {
	return r.db.Save(category).Error
}

func (r *categoryRepository) Delete(id uint) error {
	return r.db.Delete(&models.Category{}, id).Error
}
