package request

type CreateAirportRequest struct {
	Code string `json:"code" binding:"required"`
	Name string `json:"name" binding:"required"`
	City string `json:"city" binding:"required"`
}

type UpdateAirportRequest struct {
	Code *string `json:"code,omitempty"`
	Name *string `json:"name,omitempty"`
	City *string `json:"city,omitempty"`
}