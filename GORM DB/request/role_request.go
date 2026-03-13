package request

type RolePost struct {
	Name string `json:"name" binding:"required"`
}

type RolePut struct {
	Name *string `json:"name" binding:"required"`
}