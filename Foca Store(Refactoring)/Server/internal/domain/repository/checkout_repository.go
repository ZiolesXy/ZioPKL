package repository

import (
	"voca-store/internal/domain/models"
)

type CheckoutRepository interface {
	Create(checkout *models.Checkout) error
	FindByID(id uint) (models.Checkout, error)
	FindByUID(uid string) (models.Checkout, error)
	FindByMidtransOrderID(orderID string) (models.Checkout, error)
	FindAll() ([]models.Checkout, error)
	FindByUserID(userID uint) ([]models.Checkout, error)
	Update(checkout *models.Checkout, updates map[string]interface{}) error
	Save(checkout *models.Checkout) error
	Delete(checkout *models.Checkout) error
	
	// Items
	CreateItem(item *models.CheckoutItem) error
	
	// Transaction support
	WithTransaction(fn func(repo CheckoutRepository) error) error
	UpdateProductStock(productID uint, quantity int) error
	FindCouponByCode(code string) (models.Coupon, error)
	FindUserCoupon(userID uint, couponID uint) (models.UserCoupon, error)
	UpdateUserCoupon(userCoupon *models.UserCoupon, updates map[string]interface{}) error
	FindAddressByUID(uid string, userID uint) (models.Address, error)
	FindCartItems(cartID uint, itemIDs []uint) ([]models.CartItem, error)
	DeleteCartItems(cartID uint, itemIDs []uint) error
	FindCartByUserID(userID uint) (models.Cart, error)
	
	// UID generation
	GenerateUID() (string, error)
	GenerateAddressUID() (string, error)
}
