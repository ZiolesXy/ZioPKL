package product

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"voca-store/internal/domain/dto/request"
	"voca-store/internal/domain/dto/response"
	"voca-store/internal/helper"
	"voca-store/internal/service"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{productService}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	contentType := c.GetHeader("Content-Type")
	isMultipart := strings.HasPrefix(contentType, "multipart/form-data")

	var req request.CreateProductRequest
	var imageURL, imagePublicID string

	if isMultipart {
		req.Name = c.PostForm("name")
		req.Description = c.PostForm("description")
		
		priceStr := c.PostForm("price")
		price, _ := strconv.ParseFloat(priceStr, 64)
		req.Price = price

		stockStr := c.PostForm("stock")
		stock, _ := strconv.Atoi(stockStr)
		req.Stock = stock

		catIDStr := c.PostForm("category_id")
		catID, _ := strconv.ParseUint(catIDStr, 10, 32)
		req.CategoryID = uint(catID)

		// Upload image
		if file, err := c.FormFile("image"); err == nil {
			tempPath := filepath.Join("tmp", file.Filename)
			os.MkdirAll("tmp", os.ModePerm)
			c.SaveUploadedFile(file, tempPath)
			uploadRes, _ := helper.UploadFile(tempPath, "products")
			os.Remove(tempPath)
			if uploadRes != nil {
				imageURL = uploadRes.SecureURL
				imagePublicID = uploadRes.PublicID
			}
		}
	} else {
		if err := c.ShouldBindJSON(&req); err != nil {
			response.ErrorResponse(c, http.StatusBadRequest, "invalid request body")
			return
		}
		if req.ImageURL != "" {
			uploadRes, _ := helper.UploadFile(req.ImageURL, "products")
			if uploadRes != nil {
				imageURL = uploadRes.SecureURL
				imagePublicID = uploadRes.PublicID
			}
		}
	}

	productRes, err := h.productService.Create(req, imageURL, imagePublicID)
	if err != nil {
		if imagePublicID != "" {
			helper.DeleteImage(imagePublicID)
		}
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessResponse(c, "product created", productRes)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	contentType := c.GetHeader("Content-Type")
	isMultipart := strings.HasPrefix(contentType, "multipart/form-data")

	var req request.UpdateProductRequest
	var imageURL, imagePublicID string

	if isMultipart {
		name := c.PostForm("name")
		if name != "" { req.Name = &name }
		
		desc := c.PostForm("description")
		if desc != "" { req.Description = &desc }

		priceStr := c.PostForm("price")
		if priceStr != "" {
			price, _ := strconv.ParseFloat(priceStr, 64)
			req.Price = &price
		}

		stockStr := c.PostForm("stock")
		if stockStr != "" {
			stock, _ := strconv.Atoi(stockStr)
			req.Stock = &stock
		}

		catIDStr := c.PostForm("category_id")
		if catIDStr != "" {
			catID, _ := strconv.ParseUint(catIDStr, 10, 32)
			cid := uint(catID)
			req.CategoryID = &cid
		}

		if file, err := c.FormFile("image"); err == nil {
			tempPath := filepath.Join("tmp", file.Filename)
			os.MkdirAll("tmp", os.ModePerm)
			c.SaveUploadedFile(file, tempPath)
			uploadRes, _ := helper.UploadFile(tempPath, "products")
			os.Remove(tempPath)
			if uploadRes != nil {
				imageURL = uploadRes.SecureURL
				imagePublicID = uploadRes.PublicID
			}
		}
	} else {
		if err := c.ShouldBindJSON(&req); err != nil {
			response.ErrorResponse(c, http.StatusBadRequest, "invalid request body")
			return
		}
		if req.ImageURL != nil && *req.ImageURL != "" {
			uploadRes, _ := helper.UploadFile(*req.ImageURL, "products")
			if uploadRes != nil {
				imageURL = uploadRes.SecureURL
				imagePublicID = uploadRes.PublicID
			}
		}
	}

	// Delete old image if new one uploaded
	product, _ := h.productService.GetByID(uint(id))
	oldPublicID := product.ImagePublicID

	productRes, err := h.productService.Update(uint(id), req, imageURL, imagePublicID)
	if err != nil {
		if imagePublicID != "" {
			helper.DeleteImage(imagePublicID)
		}
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if imagePublicID != "" && oldPublicID != "" {
		helper.DeleteImage(oldPublicID)
	}

	response.SuccessResponse(c, "product updated", productRes)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	product, err := h.productService.GetByID(uint(id))
	if err != nil {
		response.ErrorResponse(c, http.StatusNotFound, "product not found")
		return
	}

	if product.ImagePublicID != "" {
		helper.DeleteImage(product.ImagePublicID)
	}

	if err := h.productService.Delete(uint(id)); err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "failed to delete product")
		return
	}

	response.SuccessResponse(c, "product deleted", nil)
}

func (h *ProductHandler) GetProductBySlug(c *gin.Context) {
	slug := c.Param("slug")
	product, err := h.productService.GetBySlug(slug)
	if err != nil {
		response.ErrorResponse(c, http.StatusNotFound, "product not found")
		return
	}

	response.SuccessResponse(c, "product retrieved", product)
}

func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	products, err := h.productService.GetAll()
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "failed to fetch products")
		return
	}

	response.SuccessResponse(c, "products retrieved", products)
}

func (h *ProductHandler) DeleteAllProducts(c *gin.Context) {
	products, _ := h.productService.GetProductsWithImages()
	
	deletedCount := 0
	for _, p := range products {
		if p.ImagePublicID != "" {
			helper.DeleteImage(p.ImagePublicID)
			deletedCount++
		}
	}

	if err := h.productService.DeleteAll(); err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "failed to delete products")
		return
	}

	response.SuccessResponse(c, fmt.Sprintf("all products deleted (%d images removed)", deletedCount), nil)
}

func (h *ProductHandler) DeleteAllProductImages(c *gin.Context) {
	products, _ := h.productService.GetProductsWithImages()
	
	deletedCount := 0
	for _, p := range products {
		if p.ImagePublicID != "" {
			helper.DeleteImage(p.ImagePublicID)
			deletedCount++
		}
	}

	if err := h.productService.ClearAllDatabaseImages(); err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "failed to clear image fields")
		return
	}

	response.SuccessResponse(c, fmt.Sprintf("cleared %d images", deletedCount), nil)
}
