package handlers

import (
	"net/http"
	"voca-store/helper"
	"voca-store/response"

	"github.com/gin-gonic/gin"
)

func DeleteAllCloudinaryAssets() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := helper.DeleteAllAssets()
		if err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "failed to delete all cloudinary assets")
			return
		}

		response.SuccessResponse(c, "all cloudinary assets deleted sucesdully", nil)
	}
}