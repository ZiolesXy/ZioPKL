package request

type CreateAddressRequest struct {
	Label         string `json:"label" binding:"required"`
	RecipientName string `json:"recipient_name" binding:"required"`
	Phone         string `json:"phone" binding:"required"`
	AddressLine   string `json:"address_line" binding:"required"`
	City          string `json:"city" binding:"required"`
	Province      string `json:"province" binding:"required"`
	PostalCode    string `json:"postal_code" binding:"required"`
	IsPrimary     bool   `json:"is_primary"`
}

type UpdateAddressRequest struct {
	Label         *string `json:"label"`
	RecipientName *string `json:"recipient_name"`
	Phone         *string `json:"phone"`
	AddressLine   *string `json:"address_line"`
	City          *string `json:"city"`
	Province      *string `json:"province"`
	PostalCode    *string `json:"postal_code"`
	IsPrimary     *bool   `json:"is_primary"`
}
