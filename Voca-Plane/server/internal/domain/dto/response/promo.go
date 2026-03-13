package response

import (
	"voca-plane/internal/domain/models"
)

type PromoResponse struct {
	ID       uint    `json:"id"`
	Code     string  `json:"code"`
	Discount float64 `json:"discount"`
	IsActive bool    `json:"is_active"`
}

func ToPromoResponse(p models.PromoCode) PromoResponse {
	return PromoResponse{
		ID:       p.ID,
		Code:     p.Code,
		Discount: p.Discount,
		IsActive: p.IsActive,
	}
}
