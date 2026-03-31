package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"voca-store/helper"
	"voca-store/models"
	"voca-store/request"
	"voca-store/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetProfile(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exist := c.Get("user_id")
		if !exist {
			response.ErrorResponse(c, http.StatusUnauthorized, "user not authenticated")
			return
		}

		var user models.User
		if err := db.Preload("Role").Preload("Addresses").First(&user, userID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				response.ErrorResponse(c, http.StatusNotFound, "user not found")

			} else {
				response.ErrorResponse(c, http.StatusInternalServerError, "failed to fetch user")
			}
			return
		}

		addresses := response.BuildAddressResponses(user.Addresses)

		profileResp := response.BuildUserProfileResponse(
			user.ID,
			user.Name,
			user.Email,
			user.TelephoneNumber,
			user.ProfileImageURL,
			user.Role.Name,
			addresses,
			user.CreatedAt,
			user.UpdatedAt,
		)

		response.SuccessResponse(c, "profile retrieved successfull", profileResp)
	}
}

func UpdateProfile(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			response.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
			return
		}

		var user models.User
		if err := db.First(&user, userID).Error; err != nil {
			response.ErrorResponse(c, http.StatusNotFound, "User not found")
			return
		}

		// Store old image info for cleanup
		oldProfileImagePublicID := user.ProfileImagePublicID

		// Check if it's multipart form
		contentType := c.GetHeader("Content-Type")
		isMultipart := strings.HasPrefix(contentType, "multipart/form-data")

		updates := make(map[string]interface{})

		if isMultipart {
			// Handle multipart form data
			if name := c.PostForm("name"); name != "" {
				updates["name"] = name
			}
			if description := c.PostForm("description"); description != "" {
				updates["description"] = description
			}
			if telephoneNumber := c.PostForm("telephone_number"); telephoneNumber != "" {
				updates["telephone_number"] = telephoneNumber
			}

			// Handle profile image upload
			file, err := c.FormFile("profile_image")
			if err == nil && file != nil {
				// Save file temporarily
				tempPath := "/tmp/" + file.Filename
				if err := c.SaveUploadedFile(file, tempPath); err != nil {
					response.ErrorResponse(c, http.StatusInternalServerError, "Failed to save uploaded file")
					return
				}

				// Upload to Cloudinary
				uploadResult, err := helper.UploadFile(tempPath, "user-profiles")
				if err != nil {
					os.Remove(tempPath)
					response.ErrorResponse(c, http.StatusInternalServerError, "failed to upload profile image")
					return
				}

				updates["profile_image_url"] = uploadResult.SecureURL
				updates["profile_image_public_id"] = uploadResult.PublicID

				// Clean up temp file
				os.Remove(tempPath)
			} else if c.PostForm("profile_image_url") != "" {
				// Upload from URL
				uploadResult, err := helper.UploadFile(c.PostForm("profile_image_url"), "user-profiles")
				if err != nil {
					response.ErrorResponse(c, http.StatusInternalServerError, "Failed to upload profile image from URL")
					return
				}
				updates["profile_image_url"] = uploadResult.SecureURL
				updates["profile_image_public_id"] = uploadResult.PublicID
			}
		} else {
			// Handle JSON data
			var req request.UpdateProfileRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				response.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
				return
			}

			if req.Name != nil {
				updates["name"] = *req.Name
			}
			if req.TelephoneNumber != nil {
				updates["telephone_number"] = *req.TelephoneNumber
			}
		}

		if len(updates) == 0 {
			response.ErrorResponse(c, http.StatusBadRequest, "No fields to update")
			return
		}

		// Update user
		if err := db.Model(&user).Updates(updates).Error; err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "Failed to update profile")
			return
		}

		// DELETE OLD PROFILE IMAGE FROM CLOUDINARY IF IT WAS UPDATED
		if oldProfileImagePublicID != "" && updates["profile_image_public_id"] != nil {
			// Delete the old image since we have a new one
			if err := helper.DeleteImage(oldProfileImagePublicID); err != nil {
				fmt.Printf("Warning: Failed to delete old profile image %s: %v\n", oldProfileImagePublicID, err)
			}
		} else if oldProfileImagePublicID != "" && isMultipart && c.PostForm("profile_image_url") == "" {
			// Profile image was explicitly cleared
			if err := helper.DeleteImage(oldProfileImagePublicID); err != nil {
				fmt.Printf("Warning: Failed to delete old profile image %s: %v\n", oldProfileImagePublicID, err)
			}
		}

		// Reload user
		if err := db.Preload("Role").Preload("Addresses").First(&user, user.ID).Error; err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "Failed to reload user")
			return
		}

		addresses := response.BuildAddressResponses(user.Addresses)

		// Build profile response
		profileResp := response.BuildUserProfileResponse(
			user.ID,
			user.Name,
			user.Email,
			user.TelephoneNumber,
			user.ProfileImageURL,
			user.Role.Name,
			addresses,
			user.CreatedAt,
			user.UpdatedAt,
		)

		response.SuccessResponse(c, "Profile updated successfully", profileResp)
	}
}
