package service

import (
	"errors"
	"fmt"
	"voca-store/internal/domain/models"
	"voca-store/internal/domain/repository"
	"voca-store/internal/domain/dto/request"
	"voca-store/internal/domain/dto/response"

	"github.com/gosimple/slug"
)

type ProductService interface {
	Create(req request.CreateProductRequest, imageURL, imagePublicID string) (response.ProductResponse, error)
	GetAll() (response.ProductListResponse, error)
	GetBySlug(slug string) (response.ProductResponse, error)
	Update(id uint, req request.UpdateProductRequest, imageURL, imagePublicID string) (response.ProductResponse, error)
	Delete(id uint) error
	DeleteAll() error
	GetByID(id uint) (models.Product, error)
	GetProductsWithImages() ([]models.Product, error)
	ClearAllDatabaseImages() error
}

type productService struct {
	productRepo repository.ProductRepository
	categoryRepo repository.CategoryRepository
}

func NewProductService(productRepo repository.ProductRepository, categoryRepo repository.CategoryRepository) ProductService {
	return &productService{productRepo, categoryRepo}
}

func (s *productService) Create(req request.CreateProductRequest, imageURL, imagePublicID string) (response.ProductResponse, error) {
	// Validate category
	_, err := s.categoryRepo.FindByID(req.CategoryID)
	if err != nil {
		return response.ProductResponse{}, errors.New("category not found")
	}

	uniqueSlug, err := s.generateUniqueSlug(req.Name)
	if err != nil {
		return response.ProductResponse{}, err
	}

	product := models.Product{
		Name:          req.Name,
		Slug:          uniqueSlug,
		Description:   req.Description,
		ImageURL:      imageURL,
		ImagePublicID: imagePublicID,
		Price:         req.Price,
		Stock:         req.Stock,
		CategoryID:    req.CategoryID,
	}

	if err := s.productRepo.Create(&product); err != nil {
		return response.ProductResponse{}, err
	}

	// Reload to get category info
	fullProduct, _ := s.productRepo.FindByID(product.ID)
	return response.BuildProductResponse(fullProduct), nil
}

func (s *productService) GetAll() (response.ProductListResponse, error) {
	products, err := s.productRepo.FindAll()
	if err != nil {
		return response.ProductListResponse{}, err
	}

	return response.BuildProductListResponse(products), nil
}

func (s *productService) GetBySlug(slugStr string) (response.ProductResponse, error) {
	product, err := s.productRepo.FindBySlug(slugStr)
	if err != nil {
		return response.ProductResponse{}, err
	}

	return response.BuildProductResponse(product), nil
}

func (s *productService) Update(id uint, req request.UpdateProductRequest, imageURL, imagePublicID string) (response.ProductResponse, error) {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return response.ProductResponse{}, err
	}

	updates := make(map[string]interface{})
	if req.Name != nil {
		uniqueSlug, err := s.generateUniqueSlug(*req.Name)
		if err != nil {
			return response.ProductResponse{}, err
		}
		updates["name"] = *req.Name
		updates["slug"] = uniqueSlug
	}

	if req.Description != nil {
		updates["description"] = *req.Description
	}

	if req.Price != nil {
		updates["price"] = *req.Price
	}

	if req.Stock != nil {
		updates["stock"] = *req.Stock
	}

	if req.CategoryID != nil {
		_, err := s.categoryRepo.FindByID(*req.CategoryID)
		if err != nil {
			return response.ProductResponse{}, errors.New("category not found")
		}
		updates["category_id"] = *req.CategoryID
	}

	if imageURL != "" {
		updates["image_url"] = imageURL
		updates["image_public_id"] = imagePublicID
	}

	if err := s.productRepo.Update(&product, updates); err != nil {
		return response.ProductResponse{}, err
	}

	// Reload
	fullProduct, _ := s.productRepo.FindByID(id)
	return response.BuildProductResponse(fullProduct), nil
}

func (s *productService) Delete(id uint) error {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return err
	}

	return s.productRepo.Delete(&product)
}

func (s *productService) DeleteAll() error {
	return s.productRepo.DeleteAll()
}

func (s *productService) GetByID(id uint) (models.Product, error) {
	return s.productRepo.FindByID(id)
}

func (s *productService) GetProductsWithImages() ([]models.Product, error) {
	return s.productRepo.FindWithImages()
}

func (s *productService) ClearAllDatabaseImages() error {
	return s.productRepo.ClearAllImages()
}

func (s *productService) generateUniqueSlug(name string) (string, error) {
	baseSlug := slug.Make(name)
	uniqueSlug := baseSlug
	counter := 1

	for {
		exists, err := s.productRepo.ExistsBySlug(uniqueSlug)
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
