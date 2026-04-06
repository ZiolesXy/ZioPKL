package response

import (
	"time"
	"voca-store/internal/domain/models"
)

type CheckoutItemResponse struct {
	ID       uint                `json:"id"`
	Quantity int                 `json:"quantity"`
	Price    float64             `json:"price"`
	Product  ProductMiniResponse `json:"product"`
}

type CheckoutDetailResponse struct {
	ID             uint                     `json:"id"`
	UID            string                   `json:"uid"`
	Status         string                   `json:"status"`
	User           UserMiniResponse         `json:"user"`
	Items          []CheckoutItemResponse   `json:"items"`
	Coupon         *CouponResponse          `json:"coupon,omitempty"`
	Subtotal       float64                  `json:"subtotal"`
	DiscountAmount float64                  `json:"discount_amount"`
	TotalPrice     float64                  `json:"total_price"`
	WhatsappURL    string                   `json:"whatsapp_url"`
	SnapToken      string                   `json:"snap_token,omitempty"`
	PaymentURL     string                   `json:"payment_url,omitempty"`
	PaymentStatus  string                   `json:"payment_status"`
	Address        *CheckoutAddressResponse `json:"address,omitempty"`
	CreatedAt      time.Time                `json:"created_at"`
	UpdatedAt      time.Time                `json:"updated_at"`
}

type CheckoutListResponse struct {
	Entries []CheckoutDetailResponse `json:"entries"`
}

func BuildCheckoutDetailResponse(checkout models.Checkout) CheckoutDetailResponse {
	var items []CheckoutItemResponse

	for _, item := range checkout.Items {
		items = append(items, CheckoutItemResponse{
			ID:       item.ID,
			Quantity: item.Quantity,
			Price:    item.Price,
			Product: ProductMiniResponse{
				ID:   item.Product.ID,
				Name: item.Product.Name,
			},
		})
	}

	var coupon *CouponResponse
	if checkout.Coupon != nil {
		c := BuildCouponResponse(*checkout.Coupon)
		coupon = &c
	}

	var addressResp *CheckoutAddressResponse
	if checkout.Address != nil {
		a := BuildCheckoutAddressResponse(*checkout.Address)
		addressResp = &a
	}

	return CheckoutDetailResponse{
		ID:  checkout.ID,
		UID: checkout.UID,
		User: UserMiniResponse{
			ID:              checkout.User.ID,
			Name:            checkout.User.Name,
			Email:           checkout.User.Email,
			TelephoneNumber: checkout.User.TelephoneNumber,
		},
		Coupon:         coupon,
		Address:        addressResp,
		Subtotal:       checkout.Subtotal,
		DiscountAmount: checkout.DiscountAmount,
		TotalPrice:     checkout.TotalPrice,
		WhatsappURL:    checkout.WhatsappURL,
		SnapToken:      checkout.SnapToken,
		PaymentURL:     checkout.PaymentURL,
		PaymentStatus:  checkout.PaymentStatus,
		Status:         checkout.Status,
		Items:          items,
		CreatedAt:      checkout.CreatedAt,
		UpdatedAt:      checkout.UpdatedAt,
	}
}

func BuildCheckOutListResponse(checkouts []models.Checkout) CheckoutListResponse {
	response := []CheckoutDetailResponse{}

	for _, checkout := range checkouts {
		response = append(response, BuildCheckoutDetailResponse(checkout))
	}

	return CheckoutListResponse{
		Entries: response,
	}
}
