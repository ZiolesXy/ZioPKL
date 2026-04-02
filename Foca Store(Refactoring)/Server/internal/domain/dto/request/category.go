package request

type CreateCategoryRequest struct {
	Name string `form:"name" binding:"required"`
}

type UpdateCategoryRequest struct {
	Name *string `form:"name"`
}
