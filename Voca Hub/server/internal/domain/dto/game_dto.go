package dto

type UploadGameRequest struct {
	Title       string `form:"title" binding:"required"`
	Description string `form:"description"`
	CategoryIDs []uint `form:"category_id"`
}

type UpdateGameRequest struct {
	Title       string `form:"title"`
	Description string `form:"description"`
	CategoryIDs []uint `form:"category_id"`
}
