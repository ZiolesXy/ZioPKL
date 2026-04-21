package service

import (
	"errors"

	"server/internal/domain/models"
	domainrepo "server/internal/domain/repository"
)

type FriendService struct {
	friendRepo domainrepo.FriendRepository
	userRepo   domainrepo.UserRepository
}

func NewFriendService(friendRepo domainrepo.FriendRepository, userRepo domainrepo.UserRepository) *FriendService {
	return &FriendService{
		friendRepo: friendRepo,
		userRepo:   userRepo,
	}
}

func (s *FriendService) AddFriend(userID uint, friendID uint) error {
	if userID == friendID {
		return errors.New("cannot add yourself")
	}

	target, err := s.userRepo.FindByID(friendID)
	if err != nil {
		return err
	}
	if target == nil {
		return errors.New("friend not found")
	}

	relation, err := s.friendRepo.FindRelation(userID, friendID)
	if err != nil {
		return err
	}
	if relation != nil {
		return errors.New("friend relation already exists")
	}

	return s.friendRepo.Create(&models.Friend{
		UserID:   userID,
		FriendID: friendID,
		Status:   "pending",
	})
}

func (s *FriendService) AcceptFriend(currentUserID uint, relationID uint) error {
	relation, err := s.friendRepo.FindByID(relationID)
	if err != nil {
		return err
	}
	if relation == nil {
		return errors.New("friend request not found")
	}
	if relation.FriendID != currentUserID {
		return errors.New("not allowed")
	}
	if relation.Status != "pending" {
		return errors.New("friend request already processed")
	}
	relation.Status = "accepted"
	return s.friendRepo.Update(relation)
}

func (s *FriendService) RejectFriend(currentUserID uint, relationID uint) error {
	relation, err := s.friendRepo.FindByID(relationID)
	if err != nil {
		return err
	}
	if relation == nil {
		return errors.New("friend request not found")
	}
	if relation.FriendID != currentUserID {
		return errors.New("not allowed")
	}
	if relation.Status != "pending" {
		return errors.New("friend request already processed")
	}
	relation.Status = "rejected"
	return s.friendRepo.Update(relation)
}

func (s *FriendService) ListFriends(userID uint) ([]models.User, error) {
	relations, err := s.friendRepo.ListFriends(userID)
	if err != nil {
		return nil, err
	}

	result := make([]models.User, 0, len(relations))
	for _, relation := range relations {
		if relation.UserID == userID {
			result = append(result, relation.Friend)
			continue
		}
		result = append(result, relation.User)
	}
	return result, nil
}
