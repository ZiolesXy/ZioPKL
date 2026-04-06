package response

import (
	"time"
	"voca-store/internal/domain/models"
)

type CategoryResponse struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	Slug         string    `json:"slug"`
	IconURL      string    `json:"icon_url"`
	ProductCount int64     `json:"product_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CategoryDetailResponse struct {
	ID        uint              `json:"id"`
	Name      string            `json:"name"`
	Slug      string            `json:"slug"`
	IconURL   string            `json:"icon_url"`
	Products  []ProductResponse `json:"products"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

type CategoryListResponse struct {
	Entries []CategoryResponse `json:"entries"`
}

func BuildCategoryResponse(category models.Category, count int64) CategoryResponse {
	return CategoryResponse{
		ID:           category.ID,
		Name:         category.Name,
		Slug:         category.Slug,
		IconURL:      category.IconURL,
		ProductCount: count,
		CreatedAt:    category.CreatedAt,
		UpdatedAt:    category.UpdatedAt,
	}
}

func BuildCategoryDetailResponse(category models.Category, products []models.Product) CategoryDetailResponse {
	productResponse := BuildProductListResponse(products)

	return CategoryDetailResponse{
		ID:        category.ID,
		Name:      category.Name,
		Slug:      category.Slug,
		IconURL:   category.IconURL,
		Products:  productResponse.Entries,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}
}

func BuildCategoryListResponse(categories []models.Category, counts map[uint]int64) CategoryListResponse {
	res := []CategoryResponse{}

	for _, c := range categories {
		res = append(res, BuildCategoryResponse(c, counts[c.ID]))
	}

	return CategoryListResponse{
		Entries: res,
	}
}
