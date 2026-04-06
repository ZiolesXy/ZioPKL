package checkout

import (
	"net/http"
	"voca-store/internal/domain/dto/response"
	"voca-store/internal/service"

	"github.com/gin-gonic/gin"
)

type MidtransHandler struct {
	checkoutService service.CheckoutService
}

func NewMidtransHandler(checkoutService service.CheckoutService) *MidtransHandler {
	return &MidtransHandler{checkoutService}
}

func (h *MidtransHandler) MidtransWebhook(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "invalid payload")
		return
	}

	if err := h.checkoutService.HandleMidtransWebhook(payload); err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "webhook handled"})
}
