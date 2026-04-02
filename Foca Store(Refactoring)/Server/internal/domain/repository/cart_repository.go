package repository

import (
	"voca-store/internal/domain/models"
)

type CartRepository interface {
	FindByUserID(userID uint) (models.Cart, error)
	Create(cart *models.Cart) error
	FindItem(cartID uint, productID uint) (models.CartItem, error)
	FindItemByID(itemID uint) (models.CartItem, error)
	AddItem(item *models.CartItem) error
	UpdateItem(item *models.CartItem, updates map[string]interface{}) error
	DeleteItem(item *models.CartItem) error
	DeleteItems(cartID uint, itemIDs []uint) error
	ClearCart(cartID uint) error
	CountValidItems(cartID uint, itemIDs []uint) (int64, error)
}
