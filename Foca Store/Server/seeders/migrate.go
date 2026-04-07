package seeders

import (
	"voca-store/internal/domain/models"

	"gorm.io/gorm"
)

func MigrateAll(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Role{},
		&models.User{},
		&models.Category{},
		&models.Product{},
		&models.Cart{},
		&models.CartItem{},
		&models.Checkout{},
		&models.CheckoutItem{},
		&models.Coupon{},
		&models.UserCoupon{},
		&models.Address{},
		&models.ChatSession{},
		&models.ChatMessage{},
	)
}
