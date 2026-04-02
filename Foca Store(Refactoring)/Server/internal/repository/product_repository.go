package repository

import (
	"voca-store/internal/domain/models"
	"voca-store/internal/domain/repository"

	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) repository.ProductRepository {
	return &productRepository{db}
}

func (r *productRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) FindAll() ([]models.Product, error) {
	var products []models.Product
	err := r.db.Preload("Category").Order("id ASC").Find(&products).Error
	return products, err
}

func (r *productRepository) FindBySlug(slug string) (models.Product, error) {
	var product models.Product
	err := r.db.Preload("Category").Where("slug = ?", slug).First(&product).Error
	return product, err
}

func (r *productRepository) FindByID(id uint) (models.Product, error) {
	var product models.Product
	err := r.db.Preload("Category").First(&product, id).Error
	return product, err
}

func (r *productRepository) Update(product *models.Product, updates map[string]interface{}) error {
	return r.db.Model(product).Updates(updates).Error
}

func (r *productRepository) Delete(product *models.Product) error {
	return r.db.Delete(product).Error
}

func (r *productRepository) DeleteAll() error {
	return r.db.Where("1 = 1").Delete(&models.Product{}).Error
}

func (r *productRepository) ExistsBySlug(slug string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Product{}).Where("slug = ?", slug).Count(&count).Error
	return count > 0, err
}

func (r *productRepository) FindWithImages() ([]models.Product, error) {
	var products []models.Product
	err := r.db.Where("image_public_id != ?", "").Find(&products).Error
	return products, err
}

func (r *productRepository) ClearAllImages() error {
	return r.db.Model(&models.Product{}).Where("1 = 1").Updates(map[string]interface{}{
		"image_url":       "",
		"image_public_id": "",
	}).Error
}
