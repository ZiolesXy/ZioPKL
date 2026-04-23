package handler

import (
	"net/http"
	"server/internal/helper"
	"server/internal/service"

	"github.com/gin-gonic/gin"
)

type SystemHandler struct {
	systemService service.SystemService
}

func NewSystemHandler(systemService service.SystemService) *SystemHandler {
	return &SystemHandler{systemService}
}

func (h *SystemHandler) GetNewSecret(c *gin.Context) {
	newSecret, err := h.systemService.GenerateSecret()
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, "failed generate password")
		return
	}
	helper.Success(c, http.StatusOK, "recommendation_secret", newSecret)
}