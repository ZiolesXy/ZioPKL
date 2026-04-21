package repository

import (
	"errors"

	"gorm.io/gorm"

	"server/internal/domain/models"
	domainrepo "server/internal/domain/repository"
)

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) domainrepo.PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) Create(post *models.Post) error {
	return r.db.Create(post).Error
}

func (r *postRepository) FindByID(id uint) (*models.Post, error) {
	var post models.Post
	err := r.db.Preload("User").First(&post, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) ListAll() ([]models.Post, error) {
	var posts []models.Post
	err := r.db.Preload("User").Order("created_at desc").Find(&posts).Error
	return posts, err
}

func (r *postRepository) ListByUserID(userID uint) ([]models.Post, error) {
	var posts []models.Post
	err := r.db.Preload("User").Where("user_id = ?", userID).Order("created_at desc").Find(&posts).Error
	return posts, err
}

func (r *postRepository) Update(post *models.Post) error {
	return r.db.Save(post).Error
}

func (r *postRepository) Delete(id uint) error {
	return r.db.Delete(&models.Post{}, id).Error
}
