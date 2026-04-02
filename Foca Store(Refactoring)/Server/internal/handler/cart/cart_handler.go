package cart

import (
	"net/http"
	"strconv"
	"voca-store/internal/domain/dto/request"
	"voca-store/internal/domain/dto/response"
	"voca-store/internal/service"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	cartService service.CartService
}

func NewCartHandler(cartService service.CartService) *CartHandler {
	return &CartHandler{cartService}
}

func (h *CartHandler) ViewCart(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uint)

	cartRes, err := h.cartService.GetCart(userID)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "failed to fetch cart")
		return
	}

	response.SuccessResponse(c, "cart retrieved successfully", cartRes)
}

func (h *CartHandler) AddToCart(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uint)

	var req request.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.cartService.AddToCart(userID, req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessResponse(c, "item added to cart successfully", nil)
}

func (h *CartHandler) RemoveCartItem(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uint)

	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	if err := h.cartService.RemoveItem(userID, uint(id)); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessResponse(c, "cart item removed", nil)
}

func (h *CartHandler) RemoveCartItemMany(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uint)

	var req request.RemoveCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.cartService.RemoveItems(userID, req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessResponse(c, "cart items removed", nil)
}

func (h *CartHandler) ClearCart(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uint)

	if err := h.cartService.ClearCart(userID); err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "failed to clear cart")
		return
	}

	response.SuccessResponse(c, "cart cleared", nil)
}
