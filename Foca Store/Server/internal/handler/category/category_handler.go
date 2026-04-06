package category

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"voca-store/internal/domain/dto/request"
	"voca-store/internal/domain/dto/response"
	"voca-store/internal/helper"
	"voca-store/internal/service"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categoryService service.CategoryService
}

func NewCategoryHandler(categoryService service.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService}
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req request.CreateCategoryRequest
	if err := c.ShouldBind(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "invalid request")
		return
	}

	file, err := c.FormFile("icon")
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "icon is required")
		return
	}

	tempPath := filepath.Join("tmp", file.Filename)
	if err := os.MkdirAll("tmp", os.ModePerm); err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "failed to create temp dir")
		return
	}
	if err := c.SaveUploadedFile(file, tempPath); err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "failed to save file")
		return
	}

	uploadRes, err := helper.UploadFile(tempPath, "categories/icons")
	os.Remove(tempPath)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "failed to upload icon")
		return
	}

	categoryRes, err := h.categoryService.Create(req, uploadRes.SecureURL, uploadRes.PublicID)
	if err != nil {
		helper.DeleteImage(uploadRes.PublicID)
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessResponse(c, "category created", categoryRes)
}

func (h *CategoryHandler) GetAllCategory(c *gin.Context) {
	categories, err := h.categoryService.GetAll()
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "failed to get categories")
		return
	}

	response.SuccessResponse(c, "categories retrieved", categories)
}

func (h *CategoryHandler) GetCategoryBySlug(c *gin.Context) {
	slug := c.Param("slug")
	category, err := h.categoryService.GetBySlug(slug)
	if err != nil {
		response.ErrorResponse(c, http.StatusNotFound, "category not found")
		return
	}

	response.SuccessResponse(c, "category retrieved", category)
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "invalid ID")
		return
	}

	var req request.UpdateCategoryRequest
	if err := c.ShouldBind(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "invalid request")
		return
	}

	var iconURL, iconPublicID string
	file, err := c.FormFile("icon")
	if err == nil {
		tempPath := filepath.Join("tmp", file.Filename)
		os.MkdirAll("tmp", os.ModePerm)
		c.SaveUploadedFile(file, tempPath)
		uploadRes, _ := helper.UploadFile(tempPath, "categories/icons")
		os.Remove(tempPath)
		if uploadRes != nil {
			iconURL = uploadRes.SecureURL
			iconPublicID = uploadRes.PublicID

			// Delete old icon
			category, _ := h.categoryService.GetByID(uint(id))
			if category.IconPublicID != "" {
				helper.DeleteImage(category.IconPublicID)
			}
		}
	}

	categoryRes, err := h.categoryService.Update(uint(id), req, iconURL, iconPublicID)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessResponse(c, "category updated", categoryRes)
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "invalid ID")
		return
	}

	category, err := h.categoryService.GetByID(uint(id))
	if err != nil {
		response.ErrorResponse(c, http.StatusNotFound, "category not found")
		return
	}

	if category.IconPublicID != "" {
		helper.DeleteImage(category.IconPublicID)
	}

	if err := h.categoryService.Delete(uint(id)); err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "failed to delete category")
		return
	}

	response.SuccessResponse(c, "category deleted", nil)
}
