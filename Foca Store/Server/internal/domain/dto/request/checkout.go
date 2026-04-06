package request

type CheckoutRequest struct {
	CartItemIDs []uint  `json:"cart_item_ids" binding:"required,min=1"`
	CouponCode  *string `json:"coupon_code"`
	AddressUID  string  `json:"address_uid" binding:"required"`
}

type UpdateCheckoutStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending success failed"`
}
