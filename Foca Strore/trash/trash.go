package trash

// import (
// 	"errors"
// 	"net/http"
// 	"voca-store/models"
// 	"voca-store/request"
// 	"voca-store/response"

// 	"github.com/gin-gonic/gin"
// 	"gorm.io/gorm"
// )

// func Checkout(db *gorm.DB) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		userID, exists := c.Get("user_id")
// 		if !exists {
// 			response.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
// 			return
// 		}

// 		// Get user's cart with items
// 		var cart models.Cart
// 		if err := db.Preload("Items.Product").Where("user_id = ?", userID).First(&cart).Error; err != nil {
// 			if errors.Is(err, gorm.ErrRecordNotFound) {
// 				response.ErrorResponse(c, http.StatusNotFound, "Cart not found")
// 				return
// 			}
// 			response.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch cart")
// 			return
// 		}

// 		if len(cart.Items) == 0 {
// 			response.ErrorResponse(c, http.StatusBadRequest, "Cart is empty")
// 			return
// 		}

// 		// Calculate total price and check stock
// 		totalPrice := 0.0
// 		for _, item := range cart.Items {
// 			if item.Product == nil {
// 				response.ErrorResponse(c, http.StatusBadRequest, "Product not found in cart item")
// 				return
// 			}
// 			if item.Product.Stock < item.Quantity {
// 				response.ErrorResponse(c, http.StatusBadRequest,
// 					"Insufficient stock for product: "+item.Product.Name)
// 				return
// 			}
// 			totalPrice += item.Product.Price * float64(item.Quantity)
// 		}

// 		// Create checkout record
// 		checkout := models.Checkout{
// 			UserID:     userID.(uint),
// 			TotalPrice: totalPrice,
// 			Status:     "pending",
// 		}

// 		if err := db.Create(&checkout).Error; err != nil {
// 			response.ErrorResponse(c, http.StatusInternalServerError, "Failed to create checkout")
// 			return
// 		}

// 		// Clear cart items
// 		if err := db.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error; err != nil {
// 			// Rollback checkout if failed to clear cart
// 			db.Delete(&checkout)
// 			response.ErrorResponse(c, http.StatusInternalServerError, "Failed to clear cart")
// 			return
// 		}

// 		// Build checkout response
// 		checkoutResp := response.BuildCheckoutResponse(
// 			checkout.ID,
// 			checkout.UserID,
// 			"",
// 			checkout.TotalPrice,
// 			checkout.Status,
// 			checkout.CreatedAt,
// 			checkout.UpdatedAt,
// 		)

// 		response.SuccessResponse(c, "Checkout successful", checkoutResp)
// 	}
// }

// func UpdateCheckoutStatus(db *gorm.DB) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var req request.UpdateCheckoutStatusRequest
// 		if err := c.ShouldBindJSON(&req); err != nil {
// 			response.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
// 			return
// 		}

// 		checkoutID := c.Param("id")
// 		var checkout models.Checkout
// 		if err := db.First(&checkout, checkoutID).Error; err != nil {
// 			if errors.Is(err, gorm.ErrRecordNotFound) {
// 				response.ErrorResponse(c, http.StatusNotFound, "Checkout not found")
// 			} else {
// 				response.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch checkout")
// 			}
// 			return
// 		}

// 		if err := db.Model(&checkout).Update("status", req.Status).Error; err != nil {
// 			response.ErrorResponse(c, http.StatusInternalServerError, "Failed to update checkout status")
// 			return
// 		}

// 		// Reload checkout
// 		if err := db.First(&checkout, checkout.ID).Error; err != nil {
// 			response.ErrorResponse(c, http.StatusInternalServerError, "Failed to reload checkout")
// 			return
// 		}

// 		// Build checkout response
// 		checkoutResp := response.BuildCheckoutResponse(
// 			checkout.ID,
// 			checkout.UserID,
// 			"",
// 			checkout.TotalPrice,
// 			checkout.Status,
// 			checkout.CreatedAt,
// 			checkout.UpdatedAt,
// 		)

// 		response.SuccessResponse(c, "Checkout status updated successfully", checkoutResp)
// 	}
// }

// func BuildCategoryListResponse(categories []models.Category) CategoryListResponse {
// 	response := []CategoryResponse{}

// 	for _, c := range categories {
// 		response = append(response, BuildCategoryResponse(c))
// 	}

// 	return CategoryListResponse{
// 		Entries: response,
// 	}
// }

// func Checkout(db *gorm.DB) gin.HandlerFunc {
// 	return func(c *gin.Context) {

// 		// AUTH
// 		userIDRaw, exist := c.Get("user_id")
// 		if !exist {
// 			response.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
// 			return
// 		}

// 		userID := userIDRaw.(uint)

// 		// REQUEST
// 		var req request.CheckoutRequest

// 		if err := c.ShouldBindJSON(&req); err != nil {
// 			response.ErrorResponse(c, http.StatusBadRequest, err.Error())
// 			return
// 		}

// 		// TRANSACTION START
// 		tx := db.Begin()
// 		defer func() {
// 			if r := recover(); r != nil {
// 				tx.Rollback()
// 			}
// 		}()

// 		// GET CART
// 		var cart models.Cart

// 		if err := tx.
// 			Where("user_id = ?", userID).
// 			First(&cart).Error; err != nil {
// 			tx.Rollback()
// 			response.ErrorResponse(c, http.StatusNotFound, "cart not found")
// 			return
// 		}

// 		// GET CART ITEMS + LOCK
// 		var items []models.CartItem

// 		if err := tx.
// 			Clauses(clause.Locking{Strength: "UPDATE"}).
// 			Preload("Product").
// 			Where("cart_id = ? AND id IN ?", cart.ID, req.CartItemIDs).
// 			Find(&items).Error; err != nil {
// 			tx.Rollback()
// 			response.ErrorResponse(c, http.StatusInternalServerError, "failed fetch items")
// 			return
// 		}

// 		if len(items) != len(req.CartItemIDs) {
// 			tx.Rollback()
// 			response.ErrorResponse(c, http.StatusBadRequest, "invalid cart items")
// 			return
// 		}

// 		// CALCULATE TOTAL + UPDATE STOCK
// 		totalPrice := 0.0

// 		for _, item := range items {
// 			if item.Product.Stock < item.Quantity {
// 				tx.Rollback()
// 				response.ErrorResponse(
// 					c,
// 					http.StatusBadRequest,
// 					"insufficient stock for "+item.Product.Name,
// 				)
// 				return
// 			}

// 			err := tx.Model(&models.Product{}).
// 				Where("id = ?", item.ProductID).
// 				Update("stock", gorm.Expr("stock - ?", item.Quantity)).
// 				Error

// 			if err != nil {
// 				tx.Rollback()
// 				response.ErrorResponse(c, 500, "failed update stock")
// 				return
// 			}

// 			totalPrice += item.Product.Price * float64(item.Quantity)

// 		}

// 		// APPLY COUPON
// 		var coupon models.Coupon
// 		var userCoupon models.UserCoupon

// 		if req.CouponCode != nil && *req.CouponCode != "" {

// 			err := tx.
// 				Clauses(clause.Locking{Strength: "UPDATE"}).
// 				Where("code = ?", *req.CouponCode).
// 				First(&coupon).Error

// 			if err != nil {
// 				tx.Rollback()
// 				response.ErrorResponse(c, 400, "invalid coupon")
// 				return
// 			}

// 			// VALIDATE COUPON OWNERSHIP
// 			if err := tx.
// 				Where("user_id = ? AND coupon_id = ?", userID, coupon.ID).
// 				First(&userCoupon).Error; err != nil {
// 				tx.Rollback()
// 				response.ErrorResponse(c, 400, "coupon not claimed")
// 				return
// 			}

// 			if userCoupon.UsedAt != nil {
// 				tx.Rollback()
// 				response.ErrorResponse(c, 400, "coupon already used")
// 				return
// 			}

// 			// APPLY DISCOUNT
// 			if coupon.Type == "percentage" {
// 				totalPrice -= totalPrice * coupon.Value / 100
// 			}

// 			if coupon.Type == "fixed" {
// 				totalPrice -= coupon.Value
// 			}

// 			if totalPrice < 0 {
// 				totalPrice = 0
// 			}

// 			// SAFE INCREMENT QUOTA
// 			result := tx.Model(&models.Coupon{}).
// 				Where("id = ? AND used_count < quota", coupon.ID).
// 				Update("used_count", gorm.Expr("used_count + 1"))

// 			if result.RowsAffected == 0 {
// 				tx.Rollback()
// 				response.ErrorResponse(c, 400, "coupon quota exceeded")
// 				return
// 			}

// 			// MARK COUPON AS USED
// 			now := time.Now()
// 			if err := tx.Model(&userCoupon).Update("used_at", now).Error; err != nil {
// 				tx.Rollback()
// 				response.ErrorResponse(c, 500, "failed to mark coupon as used")
// 				return
// 			}
// 		}

// 		// LOOKUP ADDRESS
// 		var address models.Address
// 		if err := tx.
// 			Where("uid = ? AND user_id = ?", req.AddressUID, userID).
// 			First(&address).Error; err != nil {
// 			tx.Rollback()
// 			response.ErrorResponse(c, http.StatusNotFound, "address not found")
// 			return
// 		}

// 		uid, err := helper.GenerateCheckoutUID(tx)
// 		if err != nil {
// 			tx.Rollback()
// 			response.ErrorResponse(c, http.StatusInternalServerError, "failed generate uid")
// 			return
// 		}

// 		// CREATE CHECKOUT
// 		checkout := models.Checkout{
// 			UID:        uid,
// 			UserID:     userID,
// 			AddressID:  &address.ID,
// 			TotalPrice: totalPrice,
// 			Status:     "pending",
// 		}

// 		if coupon.ID != 0 {
// 			checkout.CouponID = &coupon.ID
// 		}

// 		if err := tx.Create(&checkout).Error; err != nil {
// 			tx.Rollback()
// 			response.ErrorResponse(c, 500, "failed create checkout")
// 			return
// 		}

// 		// CREATE CHECKOUT ITEMS
// 		for _, item := range items {

// 			checkoutItem := models.CheckoutItem{
// 				CheckoutID: checkout.ID,
// 				ProductID:  item.ProductID,
// 				Quantity:   item.Quantity,
// 				Price:      item.Product.Price,
// 			}

// 			if err := tx.Create(&checkoutItem).Error; err != nil {
// 				tx.Rollback()
// 				response.ErrorResponse(c, 500, "failed create checkout item")
// 				return
// 			}

// 		}

// 		// DELETE CART ITEMS
// 		if err := tx.
// 			Where("id IN ?", req.CartItemIDs).
// 			Delete(&models.CartItem{}).Error; err != nil {
// 			tx.Rollback()
// 			response.ErrorResponse(c, 500, "failed delete cart")
// 			return
// 		}

// 		// COMMIT
// 		if err := tx.Commit().Error; err != nil {
// 			response.ErrorResponse(c, 500, "transaction failed")
// 			return
// 		}

// 		// RELOAD RESULT
// 		var result models.Checkout

// 		db.
// 			Preload("User").
// 			Preload("Coupon").
// 			Preload("Address").
// 			Preload("Items").
// 			Preload("Items.Product").
// 			First(&result, checkout.ID)

// 		res := response.BuildCheckoutDetailResponse(result)
// 		response.SuccessResponse(c, "checkout created", res)
// 	}
// }