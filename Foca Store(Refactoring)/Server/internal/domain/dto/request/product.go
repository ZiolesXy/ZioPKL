package request

type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	ImageURL    string  `json:"image_url"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       int     `json:"stock" binding:"required,gte=0"`
	CategoryID  uint    `json:"category_id" binding:"required"`
}

type UpdateProductRequest struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	ImageURL    *string  `json:"image_url"`
	Price       *float64 `json:"price"`
	Stock       *int     `json:"stock"`
	CategoryID  *uint    `json:"category_id"`
}

type RemoveCartItemRequest struct {
	CartItemIDs []uint `json:"cart_item_ids" binding:"required,min=1"`
}
