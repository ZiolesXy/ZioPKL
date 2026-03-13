package response

import "voca-plane/internal/domain/models"

type AirlineResponse struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Code    string `json:"code"`
	LogoURL string `json:"logo_url"`
}

func ToAirlineResponse(a models.Airline) AirlineResponse {
	return AirlineResponse{
		ID:      a.ID,
		Name:    a.Name,
		Code:    a.Code,
		LogoURL: a.LogoURL,
	}
}