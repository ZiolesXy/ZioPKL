package request

type UserPost struct {
	Name string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
	RoleID uint `json:"role_id" binding:"required"`
}

type UserPut struct {
	Name *string `json:"name"`
	Email *string `json:"email,omitempty"`
	RoleID *uint `json:"role_id,omitempty"`
}