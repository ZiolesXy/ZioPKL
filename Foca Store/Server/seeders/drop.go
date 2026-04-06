package seeders

import (
	"voca-store/internal/domain/models"

	"gorm.io/gorm"
)

func DropAllTable(db *gorm.DB) error {
	return db.Migrator().DropTable(
		&models.User{},
		&models.Role{},
		&models.Product{},
		&models.Category{},
		&models.Cart{},
		&models.CartItem{},
		&models.Checkout{},
		&models.CheckoutItem{},
		&models.Coupon{},
		&models.UserCoupon{},
		&models.Address{},
	)
}

func DropTableExceptProductsAndCategories(db *gorm.DB) error {
	return db.Migrator().DropTable(
		&models.User{},
		&models.Role{},
		&models.Cart{},
		&models.CartItem{},
		&models.Checkout{},
		&models.CheckoutItem{},
		&models.Coupon{},
		&models.UserCoupon{},
		&models.Address{},
	)
}
