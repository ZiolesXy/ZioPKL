package handlers

import (
	"net/http"
	"voca-store/models"
	"voca-store/request"
	"voca-store/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateCoupon(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req request.CreateCouponRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			response.ErrorResponse(c, http.StatusBadRequest, "invalid request body")
			return
		}

		coupon := models.Coupon {
			Code: req.Code,
			Type: req.Type,
			Value: req.Value,
			Quota: req.Quota,
			IsActive: &req.ISActive,
			ExpiresAt: &req.ExpiresAt,
		}

		if err := db.Create(&coupon).Error; err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "failed created user")
			return
		}

		res := response.BuildCouponResponse(coupon)
		response.SuccessResponse(c, "created", res)
	}
}

func GetCoupons(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		coupons := []models.Coupon{}

		if err := db.Order("id ASC").Find(&coupons).Error; err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "failed to get coupons")
			return
		}

		res := response.BuildCouponWithRemainingListResponse(coupons)
		response.SuccessListResponse(c, "coupon retrievert succesfully", res)
	}
}

func UpdateCoupon(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var coupon models.Coupon

		// cek coupon ada atau tidak
		if err := db.First(&coupon, id).Error; err != nil {
			response.ErrorResponse(c, http.StatusBadRequest, "coupon not found")
			return
		}

		var req request.UpdateCouponRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.ErrorResponse(c, http.StatusBadRequest, "invalid request body")
			return
		}

		updates := map[string]interface{}{}

		if req.Code != nil {
			updates["code"] = *req.Code
		}
		if req.Type != nil {
			updates["type"] = *req.Type
		}
		if req.Value != nil {
			updates["value"] = *req.Value
		}
		if req.Quota != nil {
			updates["quota"] = *req.Quota
		}
		if req.ISActive != nil {
			updates["is_active"] = *req.ISActive
		}
		if req.ExpiresAt != nil {
			updates["expire_at"] = *req.ExpiresAt
		}

		if len(updates) == 0 {
			response.ErrorResponse(c, http.StatusBadRequest, "no fields to update")
			return
		}

		if err := db.Model(&coupon).Updates(updates).Error; err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "failed to update coupon")
			return
		}

		// ambil data terbaru setelah update
		db.First(&coupon, id)

		res := response.BuildCouponResponse(coupon)
		response.SuccessResponse(c, "coupon successfully updated", res)
	}
}

func DeleteCoupon(db *gorm.DB) gin.HandlerFunc {

	return func(c *gin.Context) {
		id := c.Param("id")
		if err := db.Delete(&models.Coupon{},id).Error; err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "faileed to delete")
		}
		response.SuccessResponse(c,"deleted",nil)
	}
}