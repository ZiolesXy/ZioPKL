package response

import "time"

type CartItemResponse struct {
	ID        uint            `json:"id"`
	ProductID uint            `json:"product_id"`
	Product   ProductResponse `json:"product"`
	Quantity  int             `json:"quantity"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type CartResponse struct {
	ID        uint               `json:"id"`
	UserID    uint               `json:"user_id"`
	Items     []CartItemResponse `json:"items"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

func BuildCartItemResponse(id, productID uint, product ProductResponse, quantity int, createdAt, updatedAt time.Time) CartItemResponse {
	return CartItemResponse{
		ID:        id,
		ProductID: productID,
		Product:   product,
		Quantity:  quantity,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func BuildCartResponse(id, userID uint, items []CartItemResponse, createdAt, updatedAt time.Time) CartResponse {
	return CartResponse{
		ID:        id,
		UserID:    userID,
		Items:     items,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
