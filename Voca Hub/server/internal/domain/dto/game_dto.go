package dto

type UploadGameRequest struct {
	Title       string `form:"title" binding:"required"`
	Description string `form:"description"`
}
