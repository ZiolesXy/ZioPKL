package profile

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"voca-store/internal/domain/dto/request"
	"voca-store/internal/domain/dto/response"
	"voca-store/internal/helper"
	"voca-store/internal/service"

	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	authService service.AuthService
}

func NewProfileHandler(authService service.AuthService) *ProfileHandler {
	return &ProfileHandler{authService}
}

func (h *ProfileHandler) GetProfile(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uint)

	profileRes, err := h.authService.GetProfile(userID)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "failed to get profile")
		return
	}

	response.SuccessResponse(c, "profile retrieved", profileRes)
}

func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uint)

	contentType := c.GetHeader("Content-Type")
	isMultipart := strings.HasPrefix(contentType, "multipart/form-data")

	var req request.UpdateProfileRequest
	var imageURL, imagePublicID string

	if isMultipart {
		name := c.PostForm("name")
		if name != "" {
			req.Name = &name
		}

		email := c.PostForm("email")
		if email != "" {
			req.Email = &email
		}

		telp := c.PostForm("telephone_number")
		if telp != "" {
			req.TelephoneNumber = &telp
		}

		if file, err := c.FormFile("profile_image"); err == nil {
			tempPath := filepath.Join("tmp", file.Filename)
			os.MkdirAll("tmp", os.ModePerm)
			if err := c.SaveUploadedFile(file, tempPath); err != nil {
				response.ErrorResponse(c, http.StatusInternalServerError, "failed to save uploaded file")
				return
			}

			uploadRes, err := helper.UploadFile(tempPath, "profiles")
			if err != nil {
				response.ErrorResponse(c, http.StatusInternalServerError, "failed to upload image")
				return
			}
			fmt.Println("UPLOAD RESULT:", uploadRes)
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
	}

	profileRes, err := h.authService.UpdateProfile(userID, req, imageURL, imagePublicID)
	if err != nil {
		if imagePublicID != "" {
			helper.DeleteImage(imagePublicID)
		}
		response.ErrorResponse(c, http.StatusConflict, err.Error())
		return
	}

	response.SuccessResponse(c, "profile updated", profileRes)
}

func (h *ProfileHandler) ChangePassword(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uint)

	var req request.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.authService.ChangePassword(userID, req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessResponse(c, "password changed", nil)
}
