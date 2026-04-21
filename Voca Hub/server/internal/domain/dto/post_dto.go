package dto

type CreatePostRequest struct {
	Content string `json:"content" binding:"required"`
}

type UpdatePostRequest struct {
	Content string `json:"content" binding:"required"`
}
