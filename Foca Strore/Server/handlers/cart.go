package handlers

import (
	"errors"
	"net/http"
	// "time"
	"voca-store/models"
	"voca-store/request"
	"voca-store/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ViewCart(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			response.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
			return
		}

		var cart models.Cart
		if err := db.Preload("Items.Product.Category").Where("user_id = ?", userID).First(&cart).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Create cart if not exists
				cart = models.Cart{UserID: userID.(uint)}
				if err := db.Create(&cart).Error; err != nil {
					response.ErrorResponse(c, http.StatusInternalServerError, "Failed to create cart")
					return
				}
				// Reload cart
				if err := db.Preload("Items.Product.Category").First(&cart, cart.ID).Error; err != nil {
					response.ErrorResponse(c, http.StatusInternalServerError, "Failed to load cart")
					return
				}
			} else {
				response.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch cart")
				return
			}
		}

		// Build cart items response
		var cartItemResponses []response.CartItemResponse
		for _, item := range cart.Items {
			if item.Product == nil {
				continue
			}
			productResp := response.BuildProductResponse(
				*item.Product,
			)
			cartItemResp := response.BuildCartItemResponse(
				item.ID,
				item.ProductID,
				productResp,
				item.Quantity,
				item.CreatedAt,
				item.UpdatedAt,
			)
			cartItemResponses = append(cartItemResponses, cartItemResp)
		}

		// Build cart response
		cartResp := response.BuildCartResponse(
			cart.ID,
			cart.UserID,
			cartItemResponses,
			cart.CreatedAt,
			cart.UpdatedAt,
		)

		response.SuccessResponse(c, "Cart retrieved successfully", cartResp)
	}
}

func AddToCart(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			response.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
			return
		}

		var req request.AddToCartRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
			return
		}

		// Check if product exists and has enough stock
		var product models.Product
		if err := db.First(&product, req.ProductID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				response.ErrorResponse(c, http.StatusNotFound, "Product not found")
			} else {
				response.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch product")
			}
			return
		}

		if product.Stock < req.Quantity {
			response.ErrorResponse(c, http.StatusBadRequest, "Insufficient stock")
			return
		}

		// Get user's cart
		var cart models.Cart
		if err := db.Where("user_id = ?", userID).First(&cart).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Create cart if not exists
				cart = models.Cart{UserID: userID.(uint)}
				if err := db.Create(&cart).Error; err != nil {
					response.ErrorResponse(c, http.StatusInternalServerError, "Failed to create cart")
					return
				}
			} else {
				response.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch cart")
				return
			}
		}

		// Check if product already in cart
		var cartItem models.CartItem
		if err := db.Where("cart_id = ? AND product_id = ?", cart.ID, req.ProductID).First(&cartItem).Error; err == nil {
			// Update quantity
			newQuantity := cartItem.Quantity + req.Quantity
			if newQuantity > product.Stock {
				response.ErrorResponse(c, http.StatusBadRequest, "Insufficient stock")
				return
			}
			if err := db.Model(&cartItem).Update("quantity", newQuantity).Error; err != nil {
				response.ErrorResponse(c, http.StatusInternalServerError, "Failed to update cart item")
				return
			}
		} else {
			// Create new cart item
			cartItem = models.CartItem{
				CartID:    cart.ID,
				ProductID: req.ProductID,
				Quantity:  req.Quantity,
			}
			if err := db.Create(&cartItem).Error; err != nil {
				response.ErrorResponse(c, http.StatusInternalServerError, "Failed to add item to cart")
				return
			}
		}

		response.SuccessResponse(c, "Item added to cart successfully", nil)
	}
}

func RemoveCartItem(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			response.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
			return
		}

		cartItemID := c.Param("id")
		var cartItem models.CartItem
		if err := db.Where("id = ?", cartItemID).First(&cartItem).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				response.ErrorResponse(c, http.StatusNotFound, "Cart item not found")
			} else {
				response.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch cart item")
			}
			return
		}

		// Verify cart belongs to user
		var cart models.Cart
		if err := db.Where("id = ? AND user_id = ?", cartItem.CartID, userID).First(&cart).Error; err != nil {
			response.ErrorResponse(c, http.StatusForbidden, "Access denied")
			return
		}

		if err := db.Delete(&cartItem).Error; err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "Failed to remove cart item")
			return
		}

		response.SuccessResponse(c, "Cart item removed successfully", nil)
	}
}

func RemoveCartItemMany(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDRaw, exists := c.Get("user_id")
		if !exists {
			response.ErrorResponse(c, http.StatusUnauthorized, "user not authenticated")
			return
		}
		userID := userIDRaw.(uint)

		var req request.RemoveCartItemRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		var cart models.Cart
		if err := db.Where("user_id = ?", userID).First(&cart).Error; err != nil {
			response.ErrorResponse(c, http.StatusNotFound, "cart not found")
			return
		}

		var count int64
		if err := db.Model(&models.CartItem{}).Where("cart_id = ? AND id IN ?", cart.ID, req.CartItemIDs).Count(&count).Error; err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "failed to validate cart item")
			return
		}

		if count != int64(len(req.CartItemIDs)) {
			response.ErrorResponse(c, http.StatusForbidden, "some cart items are invalid")
			return
		}

		if err := db.Where("cart_id = ? AND id IN ?", cart.ID, req.CartItemIDs).Delete(&models.CartItem{}).Error; err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "failed to delete cart item")
			return
		}

		response.SuccessResponse(c, "cart items removed succesfully", nil)
	}
}

func ClearCart(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDRaw, exists := c.Get("user_id")
		if !exists {
			response.ErrorResponse(c, http.StatusUnauthorized, "user not authenticated")
			return
		}
		userID := userIDRaw.(uint)

		if err := db.Where("cart_id IN (?)", db.Model(&models.Cart{}).Select("id").Where("user_id = ?", userID)).Delete(&models.CartItem{}).Error; err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "failed to clear cart")
		}

		response.SuccessResponse(c, "cart cleared succesfully", nil)
	}
}