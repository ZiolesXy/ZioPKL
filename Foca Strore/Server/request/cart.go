package request

type AddToCartRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity int `json:"quantity" binding:"required,gte=1"`
}