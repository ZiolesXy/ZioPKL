package response

import (
	"fmt"
	"time"
	"voca-store/internal/domain/models"
)

type AddressResponse struct {
	UID           string    `json:"uid"`
	Label         string    `json:"label"`
	RecipientName string    `json:"recipient_name"`
	Phone         string    `json:"phone"`
	AddressLine   string    `json:"address_line"`
	City          string    `json:"city"`
	Province      string    `json:"province"`
	PostalCode    string    `json:"postal_code"`
	IsPrimary     bool      `json:"is_primary"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type AddressListResponse struct {
	Entries []AddressResponse `json:"entries"`
}

type CheckoutAddressResponse struct {
	AddressUID    string `json:"address_uid"`
	RecipientName string `json:"recipient_name"`
	Phone         string `json:"phone"`
	FullAddress   string `json:"full_address"`
}

func BuildAddressResponse(a models.Address) AddressResponse {
	return AddressResponse{
		UID:           a.UID,
		Label:         a.Label,
		RecipientName: a.RecipientName,
		Phone:         a.Phone,
		AddressLine:   a.AddressLine,
		City:          a.City,
		Province:      a.Province,
		PostalCode:    a.PostalCode,
		IsPrimary:     a.IsPrimary,
		CreatedAt:     a.CreatedAt,
		UpdatedAt:     a.UpdatedAt,
	}
}

func BuildAddressListResponse(addresses []models.Address) AddressListResponse {
	entries := []AddressResponse{}
	for _, a := range addresses {
		entries = append(entries, BuildAddressResponse(a))
	}
	return AddressListResponse{Entries: entries}
}

func BuildAddressResponses(address []models.Address) []AddressResponse {
	var entries []AddressResponse
	for _, a := range address {
		entries = append(entries, BuildAddressResponse(a))
	}
	return entries
}

func BuildCheckoutAddressResponse(a models.Address) CheckoutAddressResponse {
	fullAddress := fmt.Sprintf("%s, %s, %s, %s", a.AddressLine, a.City, a.Province, a.PostalCode)
	return CheckoutAddressResponse{
		AddressUID:    a.UID,
		RecipientName: a.RecipientName,
		Phone:         a.Phone,
		FullAddress:   fullAddress,
	}
}
