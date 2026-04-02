package response

import (
	"time"
	"voca-store/internal/domain/models"
)

type CategoryMiniResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type ProductMiniResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type ProductResponse struct {
	ID          uint                  `json:"id"`
	Name        string                `json:"name"`
	Slug        string                `json:"slug"`
	Description string                `json:"description,omitempty"`
	ImageURL    *string               `json:"image_url"`
	Price       float64               `json:"price"`
	Stock       int                   `json:"stock"`
	Category    *CategoryMiniResponse `json:"category,omitempty"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
}

type ProductListResponse struct {
	Entries []ProductResponse `json:"entries"`
}

func BuildProductResponse(product models.Product) ProductResponse {
	var CategoryResp *CategoryMiniResponse
	if product.Category != nil {
		CategoryResp = &CategoryMiniResponse{
			ID:   product.Category.ID,
			Name: product.Category.Name,
		}
	}

	var imageURL *string
	if product.ImageURL != "" {
		imageURL = &product.ImageURL
	}

	return ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Slug:        product.Slug,
		Description: product.Description,
		ImageURL:    imageURL,
		Price:       product.Price,
		Stock:       product.Stock,
		Category:    CategoryResp,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
}

func BuildProductListResponse(products []models.Product) ProductListResponse {
	responses := make([]ProductResponse, 0, len(products))

	for _, p := range products {
		responses = append(responses, BuildProductResponse(p))
	}

	return ProductListResponse{
		Entries: responses,
	}
}
