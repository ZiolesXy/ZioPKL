package service

import (
	"server/internal/domain/models"
	domainrepo "server/internal/domain/repository"
)

type UserService struct {
	userRepo domainrepo.UserRepository
}

func NewUserService(userRepo domainrepo.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) ListUsers() ([]models.User, error) {
	return s.userRepo.List()
}

func (s *UserService) GetByID(id uint) (*models.User, error) {
	return s.userRepo.FindByID(id)
}
