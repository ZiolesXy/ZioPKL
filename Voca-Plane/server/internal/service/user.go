package service

import (
	"context"
	"errors"
	"voca-plane/internal/domain/dto/request"
	"voca-plane/internal/domain/dto/response"
	"voca-plane/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetProfile(ctx context.Context, userID uint) (*response.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	res := response.ToUserResponse(*user)
	return &res, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, userID uint, req request.UpdateProfileRequest) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}

	if req.Name != nil {
		user.Name = *req.Name
	}

	if req.Email != nil && *req.Email != user.Email {
		_, err := s.userRepo.FindByEmail(ctx, *req.Email)
		if err == nil {
			return errors.New("email already registered")
		}
		user.Email = *req.Email
	}

	if req.Password != nil && *req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)
	}

	return s.userRepo.Update(ctx, user)
}