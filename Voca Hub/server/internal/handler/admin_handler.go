package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"server/internal/domain/dto"
	"server/internal/helper"
	"server/internal/service"
)

type AdminHandler struct {
	adminService *service.AdminService
	gameService  *service.GameService
}

func NewAdminHandler(adminService *service.AdminService, gameService *service.GameService) *AdminHandler {
	return &AdminHandler{
		adminService: adminService,
		gameService:  gameService,
	}
}

func (h *AdminHandler) Dashboard(c *gin.Context) {
	stats, err := h.adminService.Dashboard()
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.Success(c, http.StatusOK, "dashboard fetched", stats)
}

func (h *AdminHandler) ListUsers(c *gin.Context) {
	users, err := h.adminService.ListUsers()
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.Success(c, http.StatusOK, "users fetched", helper.WrapListIfNeeded(users))
}

func (h *AdminHandler) ListGames(c *gin.Context) {
	games, err := h.adminService.ListGames()
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.Success(c, http.StatusOK, "games fetched", helper.WrapListIfNeeded(dto.BuildGameResponses(games, h.gameService.BuildThumbnailURL)))
}

func (h *AdminHandler) ApproveGame(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.gameService.ApproveGame(uint(id)); err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	helper.Success(c, http.StatusOK, "game approved", nil)
}

func (h *AdminHandler) RejectGame(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.gameService.RejectGame(uint(id)); err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	helper.Success(c, http.StatusOK, "game rejected", nil)
}
