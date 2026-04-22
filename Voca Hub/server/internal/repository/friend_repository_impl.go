package repository

import (
	"errors"

	"gorm.io/gorm"

	"server/internal/domain/models"
	domainrepo "server/internal/domain/repository"
)

type friendRepository struct {
	db *gorm.DB
}

func NewFriendRepository(db *gorm.DB) domainrepo.FriendRepository {
	return &friendRepository{db: db}
}

func (r *friendRepository) FindByID(id uint) (*models.Friend, error) {
	var relation models.Friend
	err := r.db.Preload("User").Preload("Friend").First(&relation, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &relation, nil
}

func (r *friendRepository) FindRelation(userID uint, friendID uint) (*models.Friend, error) {
	var relation models.Friend
	err := r.db.Where(
		"(user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)",
		userID, friendID, friendID, userID,
	).First(&relation).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &relation, nil
}

func (r *friendRepository) Create(friend *models.Friend) error {
	return r.db.Create(friend).Error
}

func (r *friendRepository) Update(friend *models.Friend) error {
	return r.db.Save(friend).Error
}

func (r *friendRepository) ListFriends(userID uint) ([]models.Friend, error) {
	var relations []models.Friend
	err := r.db.Preload("User").Preload("Friend").
		Where("status = ? AND (user_id = ? OR friend_id = ?)", "accepted", userID, userID).
		Find(&relations).Error
	return relations, err
}

func (r *friendRepository) ListPendingRequests(friendID uint) ([]models.Friend, error) {
	var relations []models.Friend
	err := r.db.Preload("User").Preload("Friend").
		Where("status = ? AND friend_id = ?", "pending", friendID).
		Order("id desc").
		Find(&relations).Error
	return relations, err
}
