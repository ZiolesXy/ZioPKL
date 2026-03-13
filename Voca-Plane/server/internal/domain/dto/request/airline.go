package request

type CreateAirlineRequest struct {
	Name    string `json:"name" form:"name" binding:"required"`
	Code    string `json:"code" form:"code" binding:"required"`
	LogoURL string `json:"logo_url" form:"logo_url"`
}

type UpdateAirlineRequest struct {
	Name    *string `json:"name" form:"name"`
	Code    *string `json:"code" form:"code"`
	LogoURL *string `json:"logo_url" form:"logo_url"`
}