package repository

import "server/internal/domain/models"

type FriendRepository interface {
	FindByID(id uint) (*models.Friend, error)
	FindRelation(userID uint, friendID uint) (*models.Friend, error)
	Create(friend *models.Friend) error
	Update(friend *models.Friend) error
	ListFriends(userID uint) ([]models.Friend, error)
}
