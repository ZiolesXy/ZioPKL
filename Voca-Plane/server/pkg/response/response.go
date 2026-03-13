package response

import "github.com/gin-gonic/gin"

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

func Success(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, Response{
		Success: true,
		Message: message,
		Data: data,
	})
}

func SuccessWithMeta(c *gin.Context, code int, message string, data interface{}, meta interface{}) {
	c.JSON(code, Response{
		Success: true,
		Message: message,
		Data: data,
		Meta: meta,
	})
}

func Error(c *gin.Context, code int, message string) {
	c.JSON(code, Response{
		Success: false,
		Message: message,
	})
}