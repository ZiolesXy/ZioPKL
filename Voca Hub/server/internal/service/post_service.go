package service

import (
	"errors"
	"strings"

	"server/internal/domain/dto"
	"server/internal/domain/models"
	domainrepo "server/internal/domain/repository"
)

type PostService struct {
	postRepo domainrepo.PostRepository
}

func NewPostService(postRepo domainrepo.PostRepository) *PostService {
	return &PostService{postRepo: postRepo}
}

func (s *PostService) Create(userID uint, content string) (*dto.PostResponse, error) {
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
	savedPost, err := s.postRepo.FindByID(post.ID)
	if err != nil {
		return nil, err
	}
	return buildPostDTO(savedPost), nil
}

func (s *PostService) ListAll() ([]dto.PostResponse, error) {
	posts, err := s.postRepo.ListAll()
	if err != nil {
		return nil, err
	}
	return dto.BuildPostResponses(posts), nil
}

func (s *PostService) ListMine(userID uint) ([]dto.PostResponse, error) {
	posts, err := s.postRepo.ListByUserID(userID)
	if err != nil {
		return nil, err
	}
	return dto.BuildPostResponses(posts), nil
}

func (s *PostService) GetByID(id uint) (*dto.PostResponse, error) {
	post, err := s.postRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, errors.New("post not found")
	}
	return buildPostDTO(post), nil
}

func (s *PostService) Update(id uint, actor *models.User, content string) (*dto.PostResponse, error) {
	post, err := s.postRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, errors.New("post not found")
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
	savedPost, err := s.postRepo.FindByID(post.ID)
	if err != nil {
		return nil, err
	}
	return buildPostDTO(savedPost), nil
}

func (s *PostService) Delete(id uint, actor *models.User) error {
	post, err := s.postRepo.FindByID(id)
	if err != nil {
		return err
	}
	if post == nil {
		return errors.New("post not found")
	}
	if actor.Role != "ADMIN" && post.UserID != actor.ID {
		return errors.New("forbidden")
	}
	return s.postRepo.Delete(id)
}

func buildPostDTO(post *models.Post) *dto.PostResponse {
	if post == nil {
		return nil
	}

	response := dto.BuildPostResponse(*post)
	return &response
}
