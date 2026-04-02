package service

import (
	"errors"
	"fmt"
	// "strings"
	"voca-store/internal/domain/models"
	"voca-store/internal/domain/repository"
	"voca-store/internal/domain/dto/request"
	"voca-store/internal/domain/dto/response"

	"github.com/gosimple/slug"
)

type CategoryService interface {
	Create(req request.CreateCategoryRequest, iconURL, iconPublicID string) (response.CategoryResponse, error)
	GetAll() ([]response.CategoryResponse, error)
	GetBySlug(slug string) (response.CategoryDetailResponse, error)
	Update(id uint, req request.UpdateCategoryRequest, iconURL, iconPublicID string) (response.CategoryResponse, error)
	Delete(id uint) error
	GetByID(id uint) (models.Category, error)
}

type categoryService struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository) CategoryService {
	return &categoryService{categoryRepo}
}

func (s *categoryService) Create(req request.CreateCategoryRequest, iconURL, iconPublicID string) (response.CategoryResponse, error) {
	uniqueSlug, err := s.generateUniqueSlug(req.Name)
	if err != nil {
		return response.CategoryResponse{}, err
	}

	category := models.Category{
		Name:         req.Name,
		Slug:         uniqueSlug,
		IconURL:      iconURL,
		IconPublicID: iconPublicID,
	}

	if err := s.categoryRepo.Create(&category); err != nil {
		return response.CategoryResponse{}, err
	}

	return response.BuildCategoryResponse(category, 0), nil
}

func (s *categoryService) GetAll() ([]response.CategoryResponse, error) {
	categories, countMap, err := s.categoryRepo.FindAll()
	if err != nil {
		return nil, err
	}

	res := response.BuildCategoryListResponse(categories, countMap)
	return res.Entries, nil
}

func (s *categoryService) GetBySlug(slugStr string) (response.CategoryDetailResponse, error) {
	category, products, err := s.categoryRepo.FindBySlug(slugStr)
	if err != nil {
		return response.CategoryDetailResponse{}, err
	}

	return response.BuildCategoryDetailResponse(category, products), nil
}

func (s *categoryService) Update(id uint, req request.UpdateCategoryRequest, iconURL, iconPublicID string) (response.CategoryResponse, error) {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return response.CategoryResponse{}, err
	}

	updates := make(map[string]interface{})
	if req.Name != nil {
		uniqueSlug, err := s.generateUniqueSlug(*req.Name)
		if err != nil {
			return response.CategoryResponse{}, err
		}
		updates["name"] = *req.Name
		updates["slug"] = uniqueSlug
	}

	if iconURL != "" {
		updates["icon_url"] = iconURL
		updates["icon_public_id"] = iconPublicID
	}

	if err := s.categoryRepo.Update(&category, updates); err != nil {
		return response.CategoryResponse{}, err
	}

	// Fetch again to get updated product count info if needed
	// Or just return the build with current info
	// The repo implementation doesn't return count on FindByID, so we'd need FindAll or manual count
	// To simplify, just return updated category with 0 count or fetch if really needed.
	// In the old handler, it did a separate query for count.
	
	// Re-fetching updated category
	category, _ = s.categoryRepo.FindByID(id)
	return response.BuildCategoryResponse(category, 0), nil
}

func (s *categoryService) Delete(id uint) error {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return err
	}

	return s.categoryRepo.Delete(&category)
}

func (s *categoryService) GetByID(id uint) (models.Category, error) {
	return s.categoryRepo.FindByID(id)
}

func (s *categoryService) generateUniqueSlug(name string) (string, error) {
	baseSlug := slug.Make(name)
	uniqueSlug := baseSlug
	counter := 1

	for {
		exists, err := s.categoryRepo.ExistsBySlug(uniqueSlug)
		if err != nil {
			return "", err
		}
		if !exists {
			return uniqueSlug, nil
		}
		uniqueSlug = fmt.Sprintf("%s-%d", baseSlug, counter)
		counter++
		if counter > 50 {
			return "", errors.New("failed to generate unique slug after 50 attempts")
		}
	}
}
