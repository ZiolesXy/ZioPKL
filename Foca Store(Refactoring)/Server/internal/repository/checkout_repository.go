package repository

import (
	"voca-store/internal/domain/models"
	"voca-store/internal/domain/repository"
	"voca-store/internal/helper"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type checkoutRepository struct {
	db *gorm.DB
}

func NewCheckoutRepository(db *gorm.DB) repository.CheckoutRepository {
	return &checkoutRepository{db}
}

func (r *checkoutRepository) Create(checkout *models.Checkout) error {
	return r.db.Omit("User", "Address", "Coupon", "Items.Product").Create(checkout).Error
}

func (r *checkoutRepository) FindByID(id uint) (models.Checkout, error) {
	var checkout models.Checkout
	err := r.db.Preload("User").Preload("Coupon").Preload("Address").Preload("Items.Product").First(&checkout, id).Error
	return checkout, err
}

func (r *checkoutRepository) FindByUID(uid string) (models.Checkout, error) {
	var checkout models.Checkout
	err := r.db.Preload("User").Preload("Coupon").Preload("Address").Preload("Items.Product").Where("uid = ?", uid).First(&checkout).Error
	return checkout, err
}

func (r *checkoutRepository) FindByMidtransOrderID(orderID string) (models.Checkout, error) {
	var checkout models.Checkout
	err := r.db.Preload("User").Preload("Coupon").Preload("Address").Preload("Items.Product").Where("midtrans_order_id = ?", orderID).First(&checkout).Error
	return checkout, err
}

func (r *checkoutRepository) FindAll() ([]models.Checkout, error) {
	var checkouts []models.Checkout
	err := r.db.Unscoped().Preload("User").Preload("Coupon").Preload("Address").Preload("Items.Product").Order("created_at DESC").Find(&checkouts).Error
	return checkouts, err
}

func (r *checkoutRepository) FindByUserID(userID uint) ([]models.Checkout, error) {
	var checkouts []models.Checkout
	err := r.db.Preload("User").Preload("Coupon").Preload("Address").Preload("Items.Product").Where("user_id = ?", userID).Order("created_at DESC").Find(&checkouts).Error
	return checkouts, err
}

func (r *checkoutRepository) Update(checkout *models.Checkout, updates map[string]interface{}) error {
	return r.db.Model(checkout).Updates(updates).Error
}

func (r *checkoutRepository) Save(checkout *models.Checkout) error {
	return r.db.Save(checkout).Error
}

func (r *checkoutRepository) Delete(checkout *models.Checkout) error {
	return r.db.Delete(checkout).Error
}

func (r *checkoutRepository) CreateItem(item *models.CheckoutItem) error {
	return r.db.Create(item).Error
}

func (r *checkoutRepository) WithTransaction(fn func(repo repository.CheckoutRepository) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return fn(NewCheckoutRepository(tx))
	})
}

func (r *checkoutRepository) UpdateProductStock(productID uint, quantity int) error {
	return r.db.Model(&models.Product{}).Where("id = ?", productID).Update("stock", gorm.Expr("stock - ?", quantity)).Error
}

func (r *checkoutRepository) FindCouponByCode(code string) (models.Coupon, error) {
	var coupon models.Coupon
	err := r.db.Clauses(clause.Locking{Strength: "UPDATE"}).Where("code = ?", code).First(&coupon).Error
	return coupon, err
}

func (r *checkoutRepository) FindUserCoupon(userID uint, couponID uint) (models.UserCoupon, error) {
	var userCoupon models.UserCoupon
	err := r.db.Where("user_id = ? AND coupon_id = ?", userID, couponID).First(&userCoupon).Error
	return userCoupon, err
}

func (r *checkoutRepository) UpdateUserCoupon(userCoupon *models.UserCoupon, updates map[string]interface{}) error {
	return r.db.Model(userCoupon).Updates(updates).Error
}

func (r *checkoutRepository) FindAddressByUID(uid string, userID uint) (models.Address, error) {
	var address models.Address
	err := r.db.Where("uid = ? AND user_id = ?", uid, userID).First(&address).Error
	return address, err
}

func (r *checkoutRepository) FindCartItems(cartID uint, itemIDs []uint) ([]models.CartItem, error) {
	var items []models.CartItem
	err := r.db.Clauses(clause.Locking{Strength: "UPDATE"}).Preload("Product").Where("cart_id = ? AND id IN ?", cartID, itemIDs).Find(&items).Error
	return items, err
}

func (r *checkoutRepository) DeleteCartItems(cartID uint, itemIDs []uint) error {
	return r.db.Where("cart_id = ? AND id IN ?", cartID, itemIDs).Delete(&models.CartItem{}).Error
}

func (r *checkoutRepository) FindCartByUserID(userID uint) (models.Cart, error) {
	var cart models.Cart
	err := r.db.Where("user_id = ?", userID).First(&cart).Error
	return cart, err
}

func (r *checkoutRepository) GenerateUID() (string, error) {
	return helper.GenerateCheckoutUID(r.db)
}

func (r *checkoutRepository) GenerateAddressUID() (string, error) {
	return helper.NewGenerateAddressUID(r.db)
}
