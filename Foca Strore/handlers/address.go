package handlers

import (
	"net/http"
	"voca-store/helper"
	"voca-store/models"
	"voca-store/request"
	"voca-store/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateAddress(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDRaw, exists := c.Get("user_id")
		if !exists {
			response.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
			return
		}
		userID := userIDRaw.(uint)

		var req request.CreateAddressRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		tx := db.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		uid, err := helper.NewGenerateAddressUID(tx)
		if err != nil {
			tx.Rollback()
			response.ErrorResponse(c, http.StatusInternalServerError, "failed to generate address UID")
			return
		}

		address := models.Address{
			UID:           uid,
			UserID:        userID,
			Label:         req.Label,
			RecipientName: req.RecipientName,
			Phone:         req.Phone,
			AddressLine:   req.AddressLine,
			City:          req.City,
			Province:      req.Province,
			PostalCode:    req.PostalCode,
			IsPrimary:     req.IsPrimary,
		}

		if address.IsPrimary {
			if err := tx.Model(&models.Address{}).Where("user_id = ?", userID).Update("is_primary", false).Error; err != nil {
				tx.Rollback()
				response.ErrorResponse(c, http.StatusInternalServerError, "failed to unset existing primary address")
				return
			}
		}

		if err := tx.Create(&address).Error; err != nil {
			tx.Rollback()
			response.ErrorResponse(c, http.StatusInternalServerError, "failed to create address")
			return
		}

		if err := tx.Commit().Error; err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "transaction failed")
			return
		}

		res := response.BuildAddressResponse(address)
		response.SuccessResponse(c, "address created", res)
	}
}

func GetMyAddresses(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDRaw, exists := c.Get("user_id")
		if !exists {
			response.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
			return
		}
		userID := userIDRaw.(uint)

		var addresses []models.Address
		if err := db.
			Where("user_id = ?", userID).
			Order("created_at DESC").
			Find(&addresses).Error; err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "failed to fetch addresses")
			return
		}

		res := response.BuildAddressListResponse(addresses)
		response.SuccessListResponse(c, "address list fetched", res)
	}
}

func GetAddressByUID(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDRaw, exists := c.Get("user_id")
		if !exists {
			response.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
			return
		}
		userID := userIDRaw.(uint)

		uid := c.Param("uid")

		var address models.Address
		if err := db.
			Where("uid = ? AND user_id = ?", uid, userID).
			First(&address).Error; err != nil {
			response.ErrorResponse(c, http.StatusNotFound, "address not found")
			return
		}

		res := response.BuildAddressResponse(address)
		response.SuccessResponse(c, "address fetched", res)
	}
}

func UpdateAddress(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDRaw, exists := c.Get("user_id")
		if !exists {
			response.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
			return
		}
		userID := userIDRaw.(uint)

		uid := c.Param("uid")

		var address models.Address
		if err := db.
			Where("uid = ? AND user_id = ?", uid, userID).
			First(&address).Error; err != nil {
			response.ErrorResponse(c, http.StatusNotFound, "address not found")
			return
		}

		var req request.UpdateAddressRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		updates := map[string]interface{}{}
		if req.Label != nil {
			updates["label"] = *req.Label
		}
		if req.RecipientName != nil {
			updates["recipient_name"] = *req.RecipientName
		}
		if req.Phone != nil {
			updates["phone"] = *req.Phone
		}
		if req.AddressLine != nil {
			updates["address_line"] = *req.AddressLine
		}
		if req.City != nil {
			updates["city"] = *req.City
		}
		if req.Province != nil {
			updates["province"] = *req.Province
		}
		if req.PostalCode != nil {
			updates["postal_code"] = *req.PostalCode
		}

		if req.IsPrimary != nil {
			updates["is_primary"] = *req.IsPrimary
		}

		if len(updates) == 0 {
			response.ErrorResponse(c, http.StatusBadRequest, "no fields to update")
			return
		}

		tx := db.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		if req.IsPrimary != nil && *req.IsPrimary {
			if err := tx.Model(&models.Address{}).Where("user_id = ? AND uid != ?", userID, uid).Update("is_primary", false).Error; err != nil {
				tx.Rollback()
				response.ErrorResponse(c, http.StatusInternalServerError, "failed to unset existing primary address")
				return
			}
		}

		if err := tx.Model(&address).Updates(updates).Error; err != nil {
			tx.Rollback()
			response.ErrorResponse(c, http.StatusInternalServerError, "failed to update address")
			return
		}

		if err := tx.Commit().Error; err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "transaction failed")
			return
		}

		// Reload to get updated values
		db.First(&address, address.ID)

		res := response.BuildAddressResponse(address)
		response.SuccessResponse(c, "address updated", res)
	}
}

func DeleteAddress(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDRaw, exists := c.Get("user_id")
		if !exists {
			response.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
			return
		}
		userID := userIDRaw.(uint)

		uid := c.Param("uid")

		var address models.Address
		if err := db.
			Where("uid = ? AND user_id = ?", uid, userID).
			First(&address).Error; err != nil {
			response.ErrorResponse(c, http.StatusNotFound, "address not found")
			return
		}

		if err := db.Delete(&address).Error; err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "failed to delete address")
			return
		}

		response.SuccessResponse(c, "address deleted", nil)
	}
}
