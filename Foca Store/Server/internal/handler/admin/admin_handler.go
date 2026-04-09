package admin

import (
	"net/http"
	"voca-store/internal/domain/dto/response"
	"voca-store/internal/service"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	service service.AdminService
}

func NewAdminHandler(service service.AdminService) *AdminHandler {
	return &AdminHandler{service}
}

func (h *AdminHandler) GetDashboard(c *gin.Context) {
	data, err := h.service.GetDashboardData()
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "failed to load dashboard")
		return
	}

	resp := response.BuildDashboardResponse(*data)
	response.SuccessResponse(c, "dashboard loaded", resp)
}