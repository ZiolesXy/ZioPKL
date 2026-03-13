package request

type CreatePromoRequest struct {
	Code     string  `json:"code" binding:"required"`
	Discount float64 `json:"discount" binding:"required,min=0,max=100"`
	IsActive bool    `json:"is_active"`
}

type UpdatePromoRequest struct {
	Code     *string  `json:"code"`
	Discount *float64 `json:"discount"`
	IsActive *bool    `json:"is_active"`
}