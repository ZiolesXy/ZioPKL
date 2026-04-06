package request

type ClaimCouponRequest struct {
	CouponCode string `json:"coupon_code" binding:"required"`
}
