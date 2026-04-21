package service

import (
	"errors"
	"strings"

	"server/internal/domain/models"
	domainrepo "server/internal/domain/repository"
)

type PostService struct {
	postRepo domainrepo.PostRepository
}

func NewPostService(postRepo domainrepo.PostRepository) *PostService {
	return &PostService{postRepo: postRepo}
}

func (s *PostService) Create(userID uint, content string) (*models.Post, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, errors.New("content is required")
	}

	post := &models.Post{
		UserID:  userID,
		Content: content,
	}
	if err := s.postRepo.Create(post); err != nil {
		return nil, err
	}
	return s.postRepo.FindByID(post.ID)
}

func (s *PostService) ListAll() ([]models.Post, error) {
	return s.postRepo.ListAll()
}

func (s *PostService) ListMine(userID uint) ([]models.Post, error) {
	return s.postRepo.ListByUserID(userID)
}

func (s *PostService) GetByID(id uint) (*models.Post, error) {
	post, err := s.postRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, errors.New("post not found")
	}
	return post, nil
}

func (s *PostService) Update(id uint, actor *models.User, content string) (*models.Post, error) {
	post, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}
	if post.UserID != actor.ID {
		return nil, errors.New("forbidden")
	}

	content = strings.TrimSpace(content)
	if content == "" {
		return nil, errors.New("content is required")
	}

	post.Content = content
	if err := s.postRepo.Update(post); err != nil {
		return nil, err
	}
	return s.postRepo.FindByID(post.ID)
}

func (s *PostService) Delete(id uint, actor *models.User) error {
	post, err := s.GetByID(id)
	if err != nil {
		return err
	}
	if actor.Role != "ADMIN" && post.UserID != actor.ID {
		return errors.New("forbidden")
	}
	return s.postRepo.Delete(id)
}
