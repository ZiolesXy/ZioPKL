package handler

import (
	"net/http"
	"voca-plane/internal/service"
	"voca-plane/pkg/response"

	"github.com/gin-gonic/gin"
)

type SystemHandler struct {
	systemService service.SystemService
}

func NewSystemHandler(systemService service.SystemService) *SystemHandler {
	return &SystemHandler{
		systemService: systemService,
	}
}

func (h *SystemHandler) Seed(c *gin.Context) {
	if err := h.systemService.ResetAndSeed(c.Request.Context()); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "database reset and seeded successfully", nil)
}

func (h *SystemHandler) Reset(c *gin.Context) {
	if err := h.systemService.Reset(c.Request.Context()); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "database success succesfully", nil)
}