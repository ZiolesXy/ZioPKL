package handlers

import (
	"net/http"
	"time"
	"voca-store/models"
	"voca-store/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func ClaimCoupon(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDRaw, exists := c.Get("user_id")
		if !exists {
			response.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
			return
		}
		userID := userIDRaw.(uint)

		id := c.Param("id")

		tx := db.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		// Lock coupon row
		var coupon models.Coupon
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&coupon, id).Error; err != nil {

			tx.Rollback()
			response.ErrorResponse(c, http.StatusNotFound, "coupon not found")
			return
		}

		// Validate active
		if coupon.IsActive != nil && !*coupon.IsActive {
			tx.Rollback()
			response.ErrorResponse(c, http.StatusBadRequest, "coupon is not active")
			return
		}

		// Validate expired
		if coupon.ExpiresAt != nil && coupon.ExpiresAt.Before(time.Now()) {
			tx.Rollback()
			response.ErrorResponse(c, http.StatusBadRequest, "coupon has expired")
			return
		}

		// Check quota
		if coupon.Quota > 0 && coupon.UsedCount >= coupon.Quota {
			tx.Rollback()
			response.ErrorResponse(c, http.StatusBadRequest, "coupon quota exceeded")
			return
		}

		// Check duplicate claim
		var existing models.UserCoupon
		err := tx.
			Where("user_id = ? AND coupon_id = ?", userID, coupon.ID).
			First(&existing).Error

		if err == nil {
			tx.Rollback()
			response.ErrorResponse(c, http.StatusConflict, "coupon already claimed")
			return
		}

		// Increment used_count
		if coupon.Quota > 0 {
			if err := tx.Model(&models.Coupon{}).
				Where("id = ?", coupon.ID).
				Update("used_count", gorm.Expr("used_count + 1")).Error; err != nil {

				tx.Rollback()
				response.ErrorResponse(c, http.StatusInternalServerError, "failed update quota")
				return
			}
		}

		// Create claim
		userCoupon := models.UserCoupon{
			UserID:   userID,
			CouponID: coupon.ID,
		}

		if err := tx.Create(&userCoupon).Error; err != nil {
			tx.Rollback()
			response.ErrorResponse(c, http.StatusInternalServerError, "failed to claim coupon")
			return
		}

		if err := tx.Commit().Error; err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "transaction failed")
			return
		}

		db.Preload("Coupon").First(&userCoupon, userCoupon.ID)

		res := response.BuildUserCouponResponse(userCoupon)
		response.SuccessResponse(c, "Coupon claimed", res)
	}
}

func GetMyCoupons(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDRaw, exists := c.Get("user_id")
		if !exists {
			response.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
			return
		}
		userID := userIDRaw.(uint)

		var userCoupons []models.UserCoupon
		if err := db.
			Preload("Coupon").
			Where("user_id = ?", userID).
			Order("created_at DESC").
			Find(&userCoupons).Error; err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "failed to fetch coupons")
			return
		}

		res := response.BuildUserCouponListResponse(userCoupons)
		response.SuccessListResponse(c, "Your coupons fetched", res)
	}
}

func RemoveCoupon(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDRaw, exists := c.Get("user_id")
		if !exists {
			response.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
			return
		}
		userID := userIDRaw.(uint)

		id := c.Param("id")

		tx := db.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		// Find user coupon
		var userCoupon models.UserCoupon
		if err := tx.
			Where("id = ? AND user_id = ?", id, userID).
			First(&userCoupon).Error; err != nil {

			tx.Rollback()
			response.ErrorResponse(c, http.StatusNotFound, "coupon not found")
			return
		}

		// Cannot remove if already used
		if userCoupon.UsedAt != nil {
			tx.Rollback()
			response.ErrorResponse(c, http.StatusBadRequest, "cannot remove used coupon")
			return
		}

		// Lock coupon row
		var coupon models.Coupon
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&coupon, userCoupon.CouponID).Error; err != nil {

			tx.Rollback()
			response.ErrorResponse(c, http.StatusInternalServerError, "failed lock coupon")
			return
		}

		// Decrement used_count
		if coupon.Quota > 0 && coupon.UsedCount > 0 {
			if err := tx.Model(&models.Coupon{}).
				Where("id = ?", coupon.ID).
				Update("used_count", gorm.Expr("used_count - 1")).Error; err != nil {

				tx.Rollback()
				response.ErrorResponse(c, http.StatusInternalServerError, "failed update quota")
				return
			}
		}

		// Delete claim
		if err := tx.Delete(&userCoupon).Error; err != nil {
			tx.Rollback()
			response.ErrorResponse(c, http.StatusInternalServerError, "failed to remove coupon")
			return
		}

		if err := tx.Commit().Error; err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "transaction failed")
			return
		}

		response.SuccessResponse(c, "Coupon removed", nil)
	}
}