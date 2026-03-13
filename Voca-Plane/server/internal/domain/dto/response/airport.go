package response

import "voca-plane/internal/domain/models"

type AirportResponse struct {
	ID   uint   `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
	City string `json:"city"`
}

func ToAirportResponse(a models.Airport) AirportResponse {
	return AirportResponse{
		ID:   a.ID,
		Code: a.Code,
		Name: a.Name,
		City: a.City,
	}
}