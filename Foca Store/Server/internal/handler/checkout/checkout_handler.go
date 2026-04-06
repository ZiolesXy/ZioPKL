package checkout

import (
	"net/http"
	"strconv"
	"voca-store/internal/domain/dto/request"
	"voca-store/internal/domain/dto/response"
	"voca-store/internal/service"

	"github.com/gin-gonic/gin"
)

type CheckoutHandler struct {
	checkoutService service.CheckoutService
}

func NewCheckoutHandler(checkoutService service.CheckoutService) *CheckoutHandler {
	return &CheckoutHandler{checkoutService}
}

func (h *CheckoutHandler) CreateCheckout(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uint)

	var req request.CheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	checkoutRes, err := h.checkoutService.Checkout(userID, req)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessResponse(c, "checkout created", checkoutRes)
}

func (h *CheckoutHandler) ApproveCheckout(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	checkoutRes, err := h.checkoutService.ApproveCheckout(uint(id))
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessResponse(c, "checkout approved", checkoutRes)
}

func (h *CheckoutHandler) RejectCheckout(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	checkoutRes, err := h.checkoutService.RejectCheckout(uint(id))
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessResponse(c, "checkout rejected", checkoutRes)
}

func (h *CheckoutHandler) GetCheckoutList(c *gin.Context) {
	checkouts, err := h.checkoutService.GetCheckouts()
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "failed to fetch checkouts")
		return
	}

	response.SuccessResponse(c, "checkout list fetched", checkouts)
}

func (h *CheckoutHandler) GetMyCheckouts(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uint)

	checkouts, err := h.checkoutService.GetMyCheckouts(userID)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "failed to fetch your checkouts")
		return
	}

	response.SuccessResponse(c, "your checkout list fetched", checkouts)
}

func (h *CheckoutHandler) GetCheckoutByUID(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uint)
	uid := c.Param("uid")

	checkoutRes, err := h.checkoutService.GetCheckoutByUID(uid, userID)
	if err != nil {
		response.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	response.SuccessResponse(c, "checkout detail fetched", checkoutRes)
}

func (h *CheckoutHandler) DeleteCheckout(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uint)
	uid := c.Param("uid")

	if err := h.checkoutService.DeleteMyCheckout(uid, userID); err != nil {
		response.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	response.SuccessResponse(c, "checkout deleted", nil)
}
