package repository

import (
	"voca-store/internal/domain/models"
	"voca-store/internal/domain/repository"

	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) repository.CategoryRepository {
	return &categoryRepository{db}
}

type categoryWithCount struct {
	models.Category
	ProductCount int64
}

func (r *categoryRepository) Create(category *models.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) FindAll() ([]models.Category, map[uint]int64, error) {
	var rows []categoryWithCount

	err := r.db.
		Table("categories").
		Select(`
			categories.*,
			COUNT(products.id) as product_count
		`).
		Joins(`
			LEFT JOIN products
			ON products.category_id = categories.id
		`).
		Group("categories.id").
		Order("categories.id ASC").
		Scan(&rows).Error

	if err != nil {
		return nil, nil, err
	}

	var categories []models.Category
	countMap := map[uint]int64{}

	for _, row := range rows {
		categories = append(categories, row.Category)
		countMap[row.ID] = row.ProductCount
	}

	return categories, countMap, nil
}

func (r *categoryRepository) FindBySlug(slug string) (models.Category, []models.Product, error) {
	var category models.Category
	if err := r.db.Where("slug = ?", slug).First(&category).Error; err != nil {
		return category, nil, err
	}

	var products []models.Product
	if err := r.db.
		Preload("Category").
		Where("category_id = ?", category.ID).
		Order("id ASC").
		Find(&products).Error; err != nil {
		return category, nil, err
	}

	return category, products, nil
}

func (r *categoryRepository) FindByID(id uint) (models.Category, error) {
	var category models.Category
	err := r.db.First(&category, id).Error
	return category, err
}

func (r *categoryRepository) Update(category *models.Category, updates map[string]interface{}) error {
	return r.db.Model(category).Updates(updates).Error
}

func (r *categoryRepository) Delete(category *models.Category) error {
	return r.db.Delete(category).Error
}

func (r *categoryRepository) ExistsBySlug(slug string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Category{}).Where("slug = ?", slug).Count(&count).Error
	return count > 0, err
}
