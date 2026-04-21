package service

import (
	"strings"

	"server/internal/domain/models"
	domainrepo "server/internal/domain/repository"
	"server/internal/helper"
)

type UserService struct {
	userRepo    domainrepo.UserRepository
	clerkClient *helper.ClerkClient
}

func NewUserService(userRepo domainrepo.UserRepository, clerkClient *helper.ClerkClient) *UserService {
	return &UserService{
		userRepo:    userRepo,
		clerkClient: clerkClient,
	}
}

func (s *UserService) SyncUser(claims helper.ClerkClaims) (*models.User, error) {
	email, err := s.resolveEmail(claims)
	if err != nil {
		return nil, err
	}
	username, err := s.resolveUsername(claims)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByClerkID(claims.ClerkID)
	if err != nil {
		return nil, err
	}
	if user != nil {
		shouldSave := false
		if strings.TrimSpace(user.Email) == "" && email != "" {
			user.Email = email
			shouldSave = true
		}
		if username != "" && (user.Username == nil || strings.TrimSpace(*user.Username) == "") {
			user.Username = &username
			shouldSave = true
		}
		if shouldSave {
			if err := s.userRepo.Save(user); err != nil {
				return nil, err
			}
		}
		return user, nil
	}

	newUser := &models.User{
		ClerkID: claims.ClerkID,
		Email:   email,
		Username: func() *string {
			if username == "" {
				return nil
			}
			return &username
		}(),
		Role: "USER",
	}
	if err := s.userRepo.Create(newUser); err != nil {
		return nil, err
	}
	return newUser, nil
}

func (s *UserService) ListUsers() ([]models.User, error) {
	return s.userRepo.List()
}

func (s *UserService) resolveEmail(claims helper.ClerkClaims) (string, error) {
	if strings.TrimSpace(claims.Email) != "" {
		return strings.TrimSpace(claims.Email), nil
	}
	if s.clerkClient == nil {
		return "", nil
	}
	return s.clerkClient.FetchPrimaryEmail(claims.ClerkID)
}

func (s *UserService) resolveUsername(claims helper.ClerkClaims) (string, error) {
	if s.clerkClient == nil {
		return "", nil
	}
	return s.clerkClient.FetchUsername(claims.ClerkID)
}
