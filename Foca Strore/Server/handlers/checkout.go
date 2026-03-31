package handlers

import (
	"net/http"
	"time"
	"voca-store/helper"
	"voca-store/models"
	"voca-store/request"
	"voca-store/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func Checkout(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		userIDRaw, exist := c.Get("user_id")
		if !exist {
			response.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
			return
		}

		userID, ok := userIDRaw.(uint)
		if !ok {
			response.ErrorResponse(c, http.StatusInternalServerError, "invalid user context")
			return
		}

		var req request.CheckoutRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			response.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		if len(req.CartItemIDs) == 0 {
			response.ErrorResponse(c, http.StatusBadRequest, "no items selected")
			return
		}

		tx := db.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		var cart models.Cart
		if err := tx.
			Where("user_id = ?", userID).
			First(&cart).Error; err != nil {
			tx.Rollback()
			response.ErrorResponse(c, http.StatusNotFound, "cart not found")
			return
		}

		var items []models.CartItem
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Preload("Product").
			Where("cart_id = ? AND id IN ?", cart.ID, req.CartItemIDs).
			Find(&items).Error; err != nil {
			tx.Rollback()
			response.ErrorResponse(c, http.StatusInternalServerError, "failed fetch items")
			return
		}

		if len(items) != len(req.CartItemIDs) {
			tx.Rollback()
			response.ErrorResponse(c, http.StatusBadRequest, "invalid cart items")
			return
		}

		var subtotalCents int64 = 0
		var totalCents int64 = 0

		for _, item := range items {

			if item.Product.Stock < item.Quantity {
				tx.Rollback()
				response.ErrorResponse(
					c,
					http.StatusBadRequest,
					"insufficient stock for "+item.Product.Name,
				)
				return
			}

			err := tx.Model(&models.Product{}).
				Where("id = ?", item.ProductID).
				Update("stock", gorm.Expr("stock - ?", item.Quantity)).
				Error

			if err != nil {
				tx.Rollback()
				response.ErrorResponse(c, 500, "failed update stock")
				return
			}

			priceCents := int64(item.Product.Price)
			subtotalCents += priceCents * int64(item.Quantity)
		}

		totalCents = subtotalCents
		var discountAmountCents int64 = 0

		var coupon models.Coupon
		var userCoupon models.UserCoupon

		if req.CouponCode != nil && *req.CouponCode != "" {

			err := tx.
				Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("code = ?", *req.CouponCode).
				First(&coupon).Error

			if err != nil {
				tx.Rollback()
				response.ErrorResponse(c, 400, "invalid coupon")
				return
			}

			if coupon.ExpiresAt != nil && coupon.ExpiresAt.Before(time.Now()) {
				tx.Rollback()
				response.ErrorResponse(c, http.StatusBadRequest, "coupon has expired")
				return
			}

			if coupon.MinimumPurchase > 0 {
				minCent := int64(coupon.MinimumPurchase)

				if totalCents < minCent {
					tx.Rollback()
					response.ErrorResponse(c, http.StatusBadRequest, "minimum purchase not reached")
					return
				}
			}

			if err := tx.
				Where("user_id = ? AND coupon_id = ?", userID, coupon.ID).
				First(&userCoupon).Error; err != nil {
				tx.Rollback()
				response.ErrorResponse(c, 400, "coupon not claimed")
				return
			}

			if userCoupon.UsedAt != nil {
				tx.Rollback()
				response.ErrorResponse(c, 400, "coupon already used")
				return
			}

			if coupon.Type == "percentage" {
				discountAmountCents = subtotalCents * int64(coupon.Value) / 100
			}

			if coupon.Type == "fixed" {
				discountAmountCents = int64(coupon.Value)
			}

			totalCents = subtotalCents - discountAmountCents

			if totalCents < 0 {
				totalCents = 0
			}

			now := time.Now()
			if err := tx.Model(&userCoupon).
				Update("used_at", now).Error; err != nil {
				tx.Rollback()
				response.ErrorResponse(c, 500, "failed to mark coupon as used")
				return
			}
		}

		var address models.Address
		if err := tx.
			Where("uid = ? AND user_id = ?", req.AddressUID, userID).
			First(&address).Error; err != nil {
			tx.Rollback()
			response.ErrorResponse(c, http.StatusNotFound, "address not found")
			return
		}

		uid, err := helper.GenerateCheckoutUID(tx)
		if err != nil {
			tx.Rollback()
			response.ErrorResponse(c, http.StatusInternalServerError, "failed generate uid")
			return
		}

		checkout := models.Checkout{
			UID:            uid,
			UserID:         userID,
			AddressID:      &address.ID,
			Subtotal:       float64(subtotalCents),
			DiscountAmount: float64(discountAmountCents),
			TotalPrice:     float64(totalCents),
			Status:         "pending",
		}

		if coupon.ID != 0 {
			checkout.CouponID = &coupon.ID
		}

		// Preload for WhatsApp URL generation
		checkout.User = &models.User{}
		if err := tx.First(checkout.User, userID).Error; err != nil {
			tx.Rollback()
			response.ErrorResponse(c, 500, "failed load user for whatsapp")
			return
		}

		// set address
		checkout.Address = &address

		// maping item list
		// var checkoutItem []models.CheckoutItem
		for _, item := range items {
			checkout.Items = append(checkout.Items, models.CheckoutItem{
				ProductID: item.ProductID,
				Product: *item.Product,
				Quantity: item.Quantity,
				Price: item.Product.Price,
			})
		}

		checkout.WhatsappURL = helper.GenerateCheckoutWhatsappURL(checkout)

		if err := tx.Create(&checkout).Error; err != nil {
			tx.Rollback()
			response.ErrorResponse(c, 500, "failed create checkout")
			return
		}

		for _, item := range items {

			checkoutItem := models.CheckoutItem{
				CheckoutID: checkout.ID,
				ProductID:  item.ProductID,
				Quantity:   item.Quantity,
				Price:      item.Product.Price,
			}

			if err := tx.Create(&checkoutItem).Error; err != nil {
				tx.Rollback()
				response.ErrorResponse(c, 500, "failed create checkout item")
				return
			}
		}

		if err := tx.
			Where("cart_id = ? AND id IN ?", cart.ID, req.CartItemIDs).
			Delete(&models.CartItem{}).Error; err != nil {
			tx.Rollback()
			response.ErrorResponse(c, 500, "failed delete cart")
			return
		}

		// midtransOrderID := checkout.UID
		if totalCents == 0 {
			checkout.PaymentStatus = "paid"
			checkout.Status = "approved"
			checkout.MidtransOrderID = "FREE-" + checkout.UID
			checkout.SnapToken = ""
			checkout.PaymentURL = ""

			if err := tx.Save(&checkout).Error; err != nil {
				tx.Rollback()
				response.ErrorResponse(c, http.StatusInternalServerError, "failed save free checkout")
				return
			}
		} else {
			midtransResp, err := helper.CreateSnapTransaction(
				checkout.UID,
				totalCents,
				discountAmountCents,
				checkout.User.Name,
				checkout.User.Email,
				checkout.User.TelephoneNumber,
				address,
				items,
			)

			if err != nil {
				tx.Rollback()
				response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
				return
			}

			checkout.MidtransOrderID = midtransResp.OrderID
			checkout.SnapToken = midtransResp.Token
			checkout.PaymentURL = midtransResp.RedirectURL
			checkout.PaymentStatus = "pending"

			if err := tx.Save(&checkout).Error; err != nil {
				tx.Rollback()
				response.ErrorResponse(c, http.StatusInternalServerError, "failed save midtrans")
				return
			}
		}

		if err := tx.Commit().Error; err != nil {
			response.ErrorResponse(c, 500, "transaction failed")
			return
		}

		var result models.Checkout
		if err := db.
			Preload("User").
			Preload("Coupon").
			Preload("Address").
			Preload("Items").
			Preload("Items.Product").
			First(&result, checkout.ID).Error; err != nil {

			response.ErrorResponse(c, 500, "failed load checkout result")
			return
		}


		res := response.BuildCheckoutDetailResponse(result)
		response.SuccessResponse(c, "checkout created", res)
	}
}

func ApproveCheckout(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		id := c.Param("id")

		var checkout models.Checkout
		if err := db.Preload("User").
			Preload("Address").
			Preload("Items").
			Preload("Items.Product").
			First(&checkout, id).Error; err != nil {

			response.ErrorResponse(c, http.StatusNotFound, "checkout not found")
			return
		}

		if checkout.Status != "pending" {
			if checkout.Status == "rejected" {
				response.ErrorResponse(c, http.StatusBadRequest, "status already rejected cannot change")
				return
			}
			response.ErrorResponse(c, http.StatusBadRequest, "invalid status")
			return
		}

		checkout.Status = "approved"

		if err := db.Save(&checkout).Error; err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "failed update status")
			return
		}

		res := response.BuildCheckoutDetailResponse(checkout)
		response.SuccessResponse(c, "Checkout approved", res)
	}
}

func RejectCheckout(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		id := c.Param("id")

		tx := db.Begin()

		var checkout models.Checkout
		if err := tx.Preload("User").
			Preload("Address").
			Preload("Items").
			Preload("Items.Product").
			Preload("Coupon").
			// Preload("UserCoupon").
			First(&checkout, id).Error; err != nil {

			tx.Rollback()
			response.ErrorResponse(c, http.StatusNotFound, "checkout not found")
			return
		}

		if checkout.Status != "pending" {
			tx.Rollback()
			response.ErrorResponse(c, http.StatusBadRequest, "invalid status")
			return
		}

		for _, item := range checkout.Items {
			if err := tx.Model(&models.Product{}).
				Where("id = ?", item.ProductID).
				Update("stock", gorm.Expr("stock + ?", item.Quantity)).Error; err != nil {

				tx.Rollback()
				response.ErrorResponse(c, http.StatusInternalServerError, "failed restore stock")
				return
			}
		}

		checkout.Status = "rejected"

		if err := tx.Save(&checkout).Error; err != nil {
			tx.Rollback()
			response.ErrorResponse(c, http.StatusInternalServerError, "failed update status")
			return
		}

		if err := tx.Commit().Error; err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "transaction failed")
			return
		}

		res := response.BuildCheckoutDetailResponse(checkout)
		response.SuccessResponse(c, "Checkout rejected", res)
	}
}

func GetCheckout(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// status := c.Query("status")

		checkouts := []models.Checkout{}

		query := db.
			Unscoped().
			Preload("User").
			Preload("Coupon").
			Preload("Address").
			Preload("Items").
			Preload("Items.Product").
			Order("created_at DESC")

		if err := query.Find(&checkouts).Error; err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "failed to fetch checkout")
			return
		}

		res := response.BuildCheckOutListResponse(checkouts)
		response.SuccessListResponse(c, "checkout list fetched", res)
	}
}

func GetMyCheckout(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		userIDRaw, exists := c.Get("user_id")
		if !exists {
			response.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
			return
		}

		userID := userIDRaw.(uint)

		var checkouts []models.Checkout

		if err := db.
			Preload("User").
			Preload("Coupon").
			Preload("Address").
			Preload("Items").
			Preload("Items.Product").
			Where("user_id = ?", userID).
			Order("created_at DESC").
			Find(&checkouts).Error; err != nil {

			response.ErrorResponse(c, http.StatusInternalServerError, "failed to fetch checkout")
			return
		}

		res := response.BuildCheckOutListResponse(checkouts)

		response.SuccessListResponse(c, "your checkout list fetched", res)
	}
}

func GetCheckoutByUID(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		userIDRaw, exists := c.Get("user_id")
		if !exists {
			response.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
			return
		}

		userID := userIDRaw.(uint)
		uid := c.Param("uid")

		var checkout models.Checkout

		if err := db.
			Preload("User").
			Preload("Coupon").
			Preload("Address").
			Preload("Items").
			Preload("Items.Product").
			Where("uid = ? AND user_id = ?", uid, userID).
			First(&checkout).Error; err != nil {

			response.ErrorResponse(c, http.StatusNotFound, "checkout not found")
			return
		}

		res := response.BuildCheckoutDetailResponse(checkout)
		response.SuccessResponse(c, "checkout detail fetched", res)
	}
}

func DeleteMyCheckout(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDRaw, exists := c.Get("user_id")
		if !exists {
			response.ErrorResponse(c, http.StatusUnauthorized, "missing authorization header")
			return
		}
		userID := userIDRaw.(uint)
		uid := c.Param("uid")

		var checkout models.Checkout
		if err := db.Where("uid = ? AND user_id = ?", uid, userID).First(&checkout).Error; err != nil {
			response.ErrorResponse(c, http.StatusNotFound, "checkout not found")
			return
		}

		//soft delete
		if err := db.Delete(&checkout).Error; err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "failed to delete checkout")
			return
		}

		response.SuccessResponse(c, "checkout deleted succesfully", nil)
	}
}
