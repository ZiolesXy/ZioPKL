package handlers

import (
	"net/http"
	"voca-store/helper"
	"voca-store/response"

	"github.com/gin-gonic/gin"
)

func GetNewSecret(c *gin.Context) {
	newSecret, err := helper.GenerateRandomSecret(32)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "failed generate password")
		return
	}

	response.SuccessResponse(c, "recomendation_secret", newSecret)
}