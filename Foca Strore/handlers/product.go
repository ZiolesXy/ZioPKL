package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"voca-store/helper"
	"voca-store/models"
	"voca-store/request"
	"voca-store/response"
)

func CreateProduct(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		contentType := c.GetHeader("Content-Type")
		isMultipart := strings.HasPrefix(contentType, "multipart/form-data")

		var (
			name, description string
			price             float64
			stock             int
			categoryID        uint
			imageURL          string
			imagePublicID     string
			err               error
		)

		// ==============================
		// MULTIPART FORM
		// ==============================
		if isMultipart {

			name = c.PostForm("name")
			description = c.PostForm("description")

			if name == "" {
				response.ErrorResponse(c, http.StatusBadRequest, "name is required")
				return
			}

			// Parse category
			categoryStr := c.PostForm("category_id")
			if categoryStr == "" {
				response.ErrorResponse(c, http.StatusBadRequest, "category_id is required")
				return
			}

			catID64, err := strconv.ParseUint(categoryStr, 10, 64)
			if err != nil {
				response.ErrorResponse(c, http.StatusBadRequest, "invalid category_id")
				return
			}
			categoryID = uint(catID64)

			// Validate category
			var category models.Category
			if err := db.First(&category, categoryID).Error; err != nil {
				response.ErrorResponse(c, http.StatusNotFound, "category not found")
				return
			}

			// Parse price
			if priceStr := c.PostForm("price"); priceStr != "" {
				price, err = strconv.ParseFloat(priceStr, 64)
				if err != nil || price < 0 {
					response.ErrorResponse(c, http.StatusBadRequest, "invalid price")
					return
				}
			}

			// Parse stock
			if stockStr := c.PostForm("stock"); stockStr != "" {
				stock, err = strconv.Atoi(stockStr)
				if err != nil || stock < 0 {
					response.ErrorResponse(c, http.StatusBadRequest, "invalid stock")
					return
				}
			}

			// Upload image (file)
			if file, err := c.FormFile("image"); err == nil && file != nil {
				tempPath := "/tmp/" + file.Filename
				if err := c.SaveUploadedFile(file, tempPath); err != nil {
					response.ErrorResponse(c, http.StatusInternalServerError, "failed to save file")
					return
				}

				uploadResult, err := helper.UploadFile(tempPath, "products")
				os.Remove(tempPath)

				if err != nil {
					response.ErrorResponse(c, http.StatusInternalServerError, "failed to upload image")
					return
				}

				imageURL = uploadResult.SecureURL
				imagePublicID = uploadResult.PublicID
			}

		} else {
			// ==============================
			// JSON
			// ==============================

			var req request.CreateProductRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				response.ErrorResponse(c, http.StatusBadRequest, "invalid request body")
				return
			}

			name = req.Name
			description = req.Description
			price = req.Price
			stock = req.Stock
			categoryID = req.CategoryID

			if name == "" {
				response.ErrorResponse(c, http.StatusBadRequest, "name is required")
				return
			}

			if price < 0 || stock < 0 {
				response.ErrorResponse(c, http.StatusBadRequest, "price and stock must be positive")
				return
			}

			var category models.Category
			if err := db.First(&category, categoryID).Error; err != nil {
				response.ErrorResponse(c, http.StatusNotFound, "category not found")
				return
			}

			// Upload image from URL
			if req.ImageURL != "" && strings.HasPrefix(req.ImageURL, "http") {
				uploadResult, err := helper.UploadFile(req.ImageURL, "products")
				if err != nil {
					response.ErrorResponse(c, http.StatusInternalServerError, "failed to upload image")
					return
				}
				imageURL = uploadResult.SecureURL
				imagePublicID = uploadResult.PublicID
			}
		}

		// Generate unique slug
		slug, err := helper.GenerateUniqueSlug(db, name)
		if err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "failed to generate slug")
			return
		}

		product := models.Product{
			Name:          name,
			Slug:          slug,
			Description:   description,
			ImageURL:      imageURL,
			ImagePublicID: imagePublicID,
			Price:         price,
			Stock:         stock,
			CategoryID:    categoryID,
		}

		if err := db.Create(&product).Error; err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "failed to create product")
			return
		}

		db.Preload("Category").First(&product, product.ID)

		productResp := response.BuildProductResponse(product)
		response.SuccessResponse(c, "product created successfully", productResp)
	}
}

func UpdateProduct(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var product models.Product
		if err := db.First(&product, c.Param("id")).Error; err != nil {
			response.ErrorResponse(c, http.StatusNotFound, "product not found")
			return
		}

		oldImagePublicID := product.ImagePublicID
		newImagePublicID := ""
		updates := make(map[string]interface{})

		contentType := c.GetHeader("Content-Type")
		isMultipart := strings.HasPrefix(contentType, "multipart/form-data")

		// ======================================================
		// MULTIPART FORM
		// ======================================================
		if isMultipart {

			// NAME
			if name := c.PostForm("name"); name != "" {
				updates["name"] = name

				slug, err := helper.GenerateUniqueSlug(db, name)
				if err != nil {
					response.ErrorResponse(c, http.StatusInternalServerError, "failed to generate slug")
					return
				}
				updates["slug"] = slug
			}

			// DESCRIPTION
			if desc := c.PostForm("description"); desc != "" {
				updates["description"] = desc
			}

			// PRICE
			if priceStr := c.PostForm("price"); priceStr != "" {
				price, err := strconv.ParseFloat(priceStr, 64)
				if err != nil || price < 0 {
					response.ErrorResponse(c, http.StatusBadRequest, "invalid price")
					return
				}
				updates["price"] = price
			}

			// STOCK
			if stockStr := c.PostForm("stock"); stockStr != "" {
				stock, err := strconv.Atoi(stockStr)
				if err != nil || stock < 0 {
					response.ErrorResponse(c, http.StatusBadRequest, "invalid stock")
					return
				}
				updates["stock"] = stock
			}

			// CATEGORY (🔥 FIX BUG DI SINI)
			if categoryStr := c.PostForm("category_id"); categoryStr != "" {
				catID64, err := strconv.ParseUint(categoryStr, 10, 64)
				if err != nil {
					response.ErrorResponse(c, http.StatusBadRequest, "invalid category_id")
					return
				}

				var category models.Category
				if err := db.First(&category, uint(catID64)).Error; err != nil {
					response.ErrorResponse(c, http.StatusNotFound, "category not found")
					return
				}

				updates["category_id"] = uint(catID64)
			}

			// IMAGE FILE
			if file, err := c.FormFile("image"); err == nil && file != nil {
				tempPath := "/tmp/" + file.Filename
				if err := c.SaveUploadedFile(file, tempPath); err != nil {
					response.ErrorResponse(c, http.StatusInternalServerError, "failed to save file")
					return
				}

				uploadResult, err := helper.UploadFile(tempPath, "products")
				os.Remove(tempPath)

				if err != nil {
					response.ErrorResponse(c, http.StatusInternalServerError, "failed upload image")
					return
				}

				updates["image_url"] = uploadResult.SecureURL
				updates["image_public_id"] = uploadResult.PublicID
				newImagePublicID = uploadResult.PublicID
			}

			// IMAGE URL
			if imageURL := c.PostForm("image_url"); imageURL != "" {
				uploadResult, err := helper.UploadFile(imageURL, "products")
				if err != nil {
					response.ErrorResponse(c, http.StatusInternalServerError, "failed upload image")
					return
				}

				updates["image_url"] = uploadResult.SecureURL
				updates["image_public_id"] = uploadResult.PublicID
				newImagePublicID = uploadResult.PublicID
			}

		} else {

			// ======================================================
			// JSON
			// ======================================================

			var req request.UpdateProductRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				response.ErrorResponse(c, http.StatusBadRequest, "invalid request body")
				return
			}

			// NAME
			if req.Name != nil {
				updates["name"] = *req.Name

				slug, err := helper.GenerateUniqueSlug(db, *req.Name)
				if err != nil {
					response.ErrorResponse(c, http.StatusInternalServerError, "failed to generate slug")
					return
				}
				updates["slug"] = slug
			}

			// DESCRIPTION
			if req.Description != nil {
				updates["description"] = *req.Description
			}

			// PRICE
			if req.Price != nil {
				if *req.Price < 0 {
					response.ErrorResponse(c, http.StatusBadRequest, "invalid price")
					return
				}
				updates["price"] = *req.Price
			}

			// STOCK
			if req.Stock != nil {
				if *req.Stock < 0 {
					response.ErrorResponse(c, http.StatusBadRequest, "invalid stock")
					return
				}
				updates["stock"] = *req.Stock
			}

			// CATEGORY (🔥 FIX BUG DI SINI)
			if req.CategoryID != nil {
				var category models.Category
				if err := db.First(&category, *req.CategoryID).Error; err != nil {
					response.ErrorResponse(c, http.StatusNotFound, "category not found")
					return
				}

				updates["category_id"] = *req.CategoryID
			}

			// IMAGE URL
			if req.ImageURL != nil && *req.ImageURL != "" {
				uploadResult, err := helper.UploadFile(*req.ImageURL, "products")
				if err != nil {
					response.ErrorResponse(c, http.StatusInternalServerError, "failed upload image")
					return
				}

				updates["image_url"] = uploadResult.SecureURL
				updates["image_public_id"] = uploadResult.PublicID
				newImagePublicID = uploadResult.PublicID
			}
		}

		if len(updates) == 0 {
			response.ErrorResponse(c, http.StatusBadRequest, "no fields to update")
			return
		}

		if err := db.Model(&product).Updates(updates).Error; err != nil {
			if newImagePublicID != "" {
				helper.DeleteImage(newImagePublicID)
			}
			response.ErrorResponse(c, http.StatusInternalServerError, "failed to update product")
			return
		}

		// Delete old image if replaced
		if oldImagePublicID != "" && newImagePublicID != "" {
			helper.DeleteImage(oldImagePublicID)
		}

		// Reload product with category
		if err := db.Preload("Category").First(&product, product.ID).Error; err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "failed to reload product")
			return
		}

		productResp := response.BuildProductResponse(product)
		response.SuccessResponse(c, "product updated successfully", productResp)
	}
}

func DeleteProduct(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var product models.Product
		if err := db.First(&product, c.Param("id")).Error; err != nil {
			response.ErrorResponse(c, http.StatusNotFound, "product not found")
			return
		}

		if product.ImagePublicID != "" {
			helper.DeleteImage(product.ImagePublicID)
		}

		if err := db.Delete(&product).Error; err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "failed to delete product")
			return
		}

		response.SuccessResponse(c, "product deleted successfully", nil)
	}
}

func GetProductBySlug(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var product models.Product
		if err := db.Preload("Category").
			Where("slug = ?", c.Param("slug")).
			First(&product).Error; err != nil {
			response.ErrorResponse(c, http.StatusNotFound, "product not found")
			return
		}

		productResp := response.BuildProductResponse(product)
		response.SuccessResponse(c, "product retrieved successfully", productResp)
	}
}

func GetAllProducts(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var products []models.Product
		if err := db.Preload("Category").
			Order("id ASC").
			Find(&products).Error; err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "failed to fetch products")
			return
		}

		productResponses := response.BuildProductListResponse(products)
		response.SuccessListResponse(c, "products retrieved successfully", productResponses)
	}
}

func DeleteAllProducts(db *gorm.DB) gin.HandlerFunc {

	return func(c *gin.Context) {

		var products []models.Product

		// ambil semua product
		if err := db.Find(&products).Error; err != nil {

			response.ErrorResponse(
				c,
				http.StatusInternalServerError,
				"failed get products",
			)

			return
		}

		// jika kosong
		if len(products) == 0 {

			response.SuccessResponse(
				c,
				"no products to delete",
				nil,
			)

			return
		}

		// delete semua image cloudinary
		deletedCount := 0

		for _, product := range products {

			if product.ImagePublicID != "" {

				err := helper.DeleteImage(product.ImagePublicID)

				if err != nil {

					fmt.Println("Cloudinary delete failed:", err)

					continue
				}

				deletedCount++

			}
		}

		// delete semua product db
		if err := db.Where("1 = 1").Delete(&models.Product{}).Error; err != nil {

			response.ErrorResponse(
				c,
				http.StatusInternalServerError,
				"failed delete products",
			)

			return
		}

		response.SuccessResponse(
			c,
			fmt.Sprintf(
				"all products deleted successfully (%d images removed)",
				deletedCount,
			),
			nil,
		)

	}
}

func DeleteAllProductImages(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var products []models.Product

        // 1. Ambil semua produk yang memiliki ImagePublicID (biar tidak ambil semua data mubazir)
        if err := db.Where("image_public_id != ?", "").Find(&products).Error; err != nil {
            response.ErrorResponse(c, http.StatusInternalServerError, "failed to fetch products")
            return
        }

        // 2. Jika tidak ada produk dengan gambar
        if len(products) == 0 {
            response.SuccessResponse(c, "no images to delete", nil)
            return
        }

        deletedCount := 0

        // 3. Loop untuk hapus image di Cloudinary
        for _, product := range products {
            // Deteksi: Jika ada PublicID, hapus di Cloudinary
            if product.ImagePublicID != "" {
                err := helper.DeleteImage(product.ImagePublicID)
                if err != nil {
                    // Jika gagal di cloudinary, kita log tapi tetap lanjut (opsional)
                    fmt.Printf("Cloudinary delete failed for ID %s: %v\n", product.ImagePublicID, err)
                    continue 
                }
                deletedCount++
            }
        }

        // 4. Update DB: Kosongkan kolom image_url dan image_public_id untuk SEMUA produk
        // Kita menggunakan .Model(&models.Product{}) dan .Where("1 = 1") untuk menyasar semua row
        updateData := map[string]interface{}{
            "image_url":       "",
            "image_public_id": "",
        }

        if err := db.Model(&models.Product{}).Where("1 = 1").Updates(updateData).Error; err != nil {
            response.ErrorResponse(c, http.StatusInternalServerError, "failed to clear image fields in database")
            return
        }

        response.SuccessResponse(
            c,
            fmt.Sprintf("Successfully cleared %d images from Cloudinary and updated database", deletedCount),
            nil,
        )
    }
}