package system

import (
	"net/http"
	"voca-store/internal/domain/dto/response"
	"voca-store/internal/service"

	"github.com/gin-gonic/gin"
)

type SystemHandler struct {
	systemService service.SystemService
}

func NewSystemHandler(systemService service.SystemService) *SystemHandler {
	return &SystemHandler{systemService}
}

func (h *SystemHandler) ResetDatabase(c *gin.Context) {
	if err := h.systemService.ResetDatabase(); err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessResponse(c, "database reset successfully", nil)
}

func (h *SystemHandler) ResetDatabaseWithProducts(c *gin.Context) {
	if err := h.systemService.ResetDatabaseWithProducts(); err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessResponse(c, "database reset (with products) successfully", nil)
}

func (h *SystemHandler) ResetDatabasePreserveCatalog(c *gin.Context) {
	if err := h.systemService.ResetDatabasePreserveCatalog(); err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessResponse(c, "database reset (preserving catalog) successfully", nil)
}

func (h *SystemHandler) Migrate(c *gin.Context) {
	if err := h.systemService.MigrateAll(); err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessResponse(c, "migration complete", nil)
}

func (h *SystemHandler) ResetRedis(c *gin.Context) {
	if err := h.systemService.ResetRedis(); err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessResponse(c, "redis successfully reset", nil)
}

func (h *SystemHandler) DeleteAllCloudinaryAssets(c *gin.Context) {
	if err := h.systemService.DeleteAllCloudinaryAssets(); err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "failed to delete all cloudinary assets")
		return
	}
	response.SuccessResponse(c, "all cloudinary assets deleted successfully", nil)
}

func (h *SystemHandler) GetNewSecret(c *gin.Context) {
	newSecret, err := h.systemService.GenerateSecret()
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "failed generate password")
		return
	}
	response.SuccessResponse(c, "recommendation_secret", newSecret)
}

func (h *SystemHandler) SeedRoles(c *gin.Context) {
	if err := h.systemService.SeedRoles(); err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to seed roles")
		return
	}
	response.SuccessResponse(c, "Roles seeded successfully", nil)
}

func (h *SystemHandler) SeedAdmin(c *gin.Context) {
	if err := h.systemService.SeedAdmin(); err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to seed admin")
		return
	}
	response.SuccessResponse(c, "Admin seeded successfully", nil)
}

func (h *SystemHandler) SeedUsers(c *gin.Context) {
	if err := h.systemService.SeedUsers(); err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to seed users")
		return
	}
	response.SuccessResponse(c, "Users seeded successfully", nil)
}

func (h *SystemHandler) SeedProducts(c *gin.Context) {
	if err := h.systemService.SeedProducts(); err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to seed products")
		return
	}
	response.SuccessResponse(c, "Products seeded successfully", nil)
}

func (h *SystemHandler) SeedProductsFromAssets(c *gin.Context) {
	if err := h.systemService.SeedProductsFromAssets(); err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to seed products from assets")
		return
	}
	response.SuccessResponse(c, "Products seeded from assets successfully", nil)
}

func (h *SystemHandler) SyncAssetProducts(c *gin.Context) {
	if err := h.systemService.SyncAssetProducts(); err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to sync asset products")
		return
	}
	response.SuccessResponse(c, "Asset products synced successfully", nil)
}

func (h *SystemHandler) SeedCoupons(c *gin.Context) {
	if err := h.systemService.SeedCoupons(); err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to seed coupons")
		return
	}
	response.SuccessResponse(c, "Coupons seeded successfully", nil)
}

func (h *SystemHandler) SeedAll(c *gin.Context) {
	if err := h.systemService.SeedAll(); err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to seed all data")
		return
	}
	response.SuccessResponse(c, "All data seeded successfully", nil)
}

func (h *SystemHandler) SeedAllWithProducts(c *gin.Context) {
	if err := h.systemService.SeedAllWithProducts(); err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to seed all data with products")
		return
	}
	response.SuccessResponse(c, "All data with products seeded successfully", nil)
}

func (h *SystemHandler) ClearChatHistory(c *gin.Context) {
	if err := h.systemService.ClearChat(); err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to clear chat history")
		return
	}
	response.SuccessResponse(c, "Chat history cleared successfully", nil)
}