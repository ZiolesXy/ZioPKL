package handlers

import (
	"net/http"

	"voca-store/response"
	"voca-store/seeders"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SeedRoleHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		if err := seeders.SeedRoles(db); err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "Failed to seed roles")
			return
		}

		response.SuccessResponse(c, "Roles seeded successfully", nil)
	}
}

func SeedAdminHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		if err := seeders.SeedAdmin(db); err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "Failed to seed admin")
			return
		}

		response.SuccessResponse(c, "Admin seeded successfully", nil)
	}
}

func SeedUsersHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		if err := seeders.SeedUsers(db); err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "Failed to seed users")
			return
		}

		response.SuccessResponse(c, "Users seeded successfully", nil)
	}
}

func SeedProductsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		if err := seeders.SeedProducts(db); err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "Failed to seed products")
			return
		}

		response.SuccessResponse(c, "Products seeded successfully", nil)
	}
}

func SeedProductsFromAssetsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		err := seeders.SeedProductsFromAssets(db)

		if err != nil {

			response.ErrorResponse(
				c,
				http.StatusInternalServerError,
				"Failed to seed products",
			)

			return
		}

		response.SuccessResponse(
			c,
			"Products seeded from assets successfully",
			nil,
		)
	}
}

func SyncAssetProductsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		err := seeders.SyncAssetProductsWithDefaultSeed(db)
		if err != nil {
			response.ErrorResponse(
				c,
				http.StatusInternalServerError,
				"Failed to sync asset products",
			)
			return
		}

		response.SuccessResponse(
			c,
			"Asset products synced successfully",
			nil,
		)
	}
}

func SeedCouponHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := seeders.SeedCoupons(db); err != nil {

			response.ErrorResponse(
				c,
				http.StatusInternalServerError,
				"Failed to seed coupons",
			)
			return
		}

		response.SuccessResponse(
			c,
			"Coupons seeded successfully",
			nil,
		)
	}
}

func SeedAllHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		if err := seeders.SeedRoles(db); err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "Failed to seed roles")
			return
		}

		if err := seeders.SeedAdmin(db); err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "Failed to seed admin")
			return
		}

		if err := seeders.SeedCategories(db); err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "Failed to seed category")
		}

		if err := seeders.SeedUsers(db); err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "Failed to seed users")
			return
		}

		if err := seeders.SeedCoupons(db); err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "Failed to seed coupons")
		}

		response.SuccessResponse(c, "All data seeded successfully", nil)
	}
}

func SeedAllWithProductnonAssetsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		if err := seeders.SeedRoles(db); err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "Failed to seed roles")
			return
		}

		if err := seeders.SeedAdmin(db); err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "Failed to seed admin")
			return
		}

		if err := seeders.SeedCategories(db); err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "Failed to seed category")
		}

		if err := seeders.SeedUsers(db); err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "Failed to seed users")
			return
		}

		if err := seeders.SeedCoupons(db); err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "Failed to seed coupons")
		}

		if err := seeders.SeedProducts(db); err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "Failed to seed products")
		}

		response.SuccessResponse(c, "All data seeded successfully", nil)
	}
}

func ResetDatabaseHandler(db *gorm.DB, rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := seeders.ResetDatabase(db, rdb); err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		response.SuccessResponse(c, "database reset succesfully", nil)
	}
}

func ResetDatabaseWithProductsHandler(db *gorm.DB, rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := seeders.ResetDatabaseWithProduct(db, rdb); err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		response.SuccessResponse(c, "database reset succesfully", nil)
	}
}

func ResetDatabasePreserveProductsAndCategoriesHandler(db *gorm.DB, rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := seeders.ResetDatabasePreserveProductsAndCategories(db, rdb); err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		response.SuccessResponse(c, "database reset (preserving catalog) succesfully", nil)
	}
}


func MigrateHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := seeders.MigrateAll(db); err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		response.SuccessResponse(c, "migration complite", nil)
	}
}

func ResetRedis(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := seeders.ResetRedis(rdb); err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		response.SuccessResponse(c, "redis succesfully reset", nil)
	}
}