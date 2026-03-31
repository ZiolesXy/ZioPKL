package response

type RoleResp struct {
	ID uint `json:"id"`
	Name string `json:"name"`
}

func BuildRoleResponse(id uint, name string) RoleResp{
	return RoleResp{
		ID: id,
		Name: name,
	}
}