package helper

import (
	"reflect"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func Success(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, APIResponse{
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, status int, message string) {
	c.JSON(status, APIResponse{
		Message: "error",
		Error:   message,
	})
}

func WrapListIfNeeded(data interface{}) interface{} {
	v := reflect.ValueOf(data)

	// Jika nil, return langsung
	if !v.IsValid() {
		return data
	}

	// Cek apakah slice atau array → bungkus sebagai entries
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		return gin.H{
			"entries": data,
		}
	}

	// Bukan slice/array → return apa adanya
	return data
}