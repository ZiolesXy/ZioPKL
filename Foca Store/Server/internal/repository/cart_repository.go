package repository

import (
	"voca-store/internal/domain/models"
	"voca-store/internal/domain/repository"

	"gorm.io/gorm"
)

type cartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) repository.CartRepository {
	return &cartRepository{db}
}

func (r *cartRepository) FindByUserID(userID uint) (models.Cart, error) {
	var cart models.Cart
	err := r.db.Preload("Items.Product.Category").Where("user_id = ?", userID).First(&cart).Error
	return cart, err
}

func (r *cartRepository) Create(cart *models.Cart) error {
	return r.db.Create(cart).Error
}

func (r *cartRepository) FindItem(cartID uint, productID uint) (models.CartItem, error) {
	var item models.CartItem
	err := r.db.Where("cart_id = ? AND product_id = ?", cartID, productID).First(&item).Error
	return item, err
}

func (r *cartRepository) FindItemByID(itemID uint) (models.CartItem, error) {
	var item models.CartItem
	err := r.db.First(&item, itemID).Error
	return item, err
}

func (r *cartRepository) AddItem(item *models.CartItem) error {
	return r.db.Create(item).Error
}

func (r *cartRepository) UpdateItem(item *models.CartItem, updates map[string]interface{}) error {
	return r.db.Model(item).Updates(updates).Error
}

func (r *cartRepository) DeleteItem(item *models.CartItem) error {
	return r.db.Delete(item).Error
}

func (r *cartRepository) DeleteItems(cartID uint, itemIDs []uint) error {
	return r.db.Where("cart_id = ? AND id IN ?", cartID, itemIDs).Delete(&models.CartItem{}).Error
}

func (r *cartRepository) ClearCart(cartID uint) error {
	return r.db.Where("cart_id = ?", cartID).Delete(&models.CartItem{}).Error
}

func (r *cartRepository) CountValidItems(cartID uint, itemIDs []uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.CartItem{}).Where("cart_id = ? AND id IN ?", cartID, itemIDs).Count(&count).Error
	return count, err
}
