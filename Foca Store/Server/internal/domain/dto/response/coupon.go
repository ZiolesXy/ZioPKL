package response

import (
	"time"
	"voca-store/internal/domain/models"
)

type CouponResponse struct {
	ID              uint       `json:"id"`
	Code            string     `json:"code"`
	Type            string     `json:"type"`
	Value           float64    `json:"value"`
	Quota           int        `json:"quota"`
	UsedCount       int        `json:"used_count"`
	MinimumPurchase float64    `json:"minimum_purchase"`
	IsActive        bool       `json:"is_active"`
	ExpiresAt       *time.Time `json:"expires_at"`
}

type CouponListResponse struct {
	Entries []CouponResponse `json:"entries"`
}

func BuildCouponResponse(c models.Coupon) CouponResponse {
	return CouponResponse{
		ID:              c.ID,
		Code:            c.Code,
		Type:            c.Type,
		Value:           c.Value,
		Quota:           c.Quota,
		UsedCount:       c.UsedCount,
		MinimumPurchase: c.MinimumPurchase,
		IsActive:        *c.IsActive,
		ExpiresAt:       c.ExpiresAt,
	}
}

func BuildCouponListResponse(coupons []models.Coupon) CouponListResponse {
	res := []CouponResponse{}

	for _, c := range coupons {
		res = append(res, BuildCouponResponse(c))
	}

	return CouponListResponse{
		Entries: res,
	}
}

type CouponWithRemainingResponse struct {
	ID              uint       `json:"id"`
	Code            string     `json:"code"`
	Type            string     `json:"type"`
	Value           float64    `json:"value"`
	Quota           int        `json:"quota"`
	UsedCount       int        `json:"used_count"`
	Remaining       int        `json:"remaining"`
	MinimumPurchase float64    `json:"minimum_purchase"`
	IsActive        bool       `json:"is_active"`
	ExpiresAt       *time.Time `json:"expires_at"`
}

type CouponWithRemainingListResponse struct {
	Entries []CouponWithRemainingResponse `json:"entries"`
}

func BuildCouponWithRemainingResponse(c models.Coupon) CouponWithRemainingResponse {
	remaining := c.Quota - c.UsedCount
	if remaining < 0 {
		remaining = 0
	}

	return CouponWithRemainingResponse{
		ID:              c.ID,
		Code:            c.Code,
		Type:            c.Type,
		Value:           c.Value,
		Quota:           c.Quota,
		UsedCount:       c.UsedCount,
		Remaining:       remaining,
		MinimumPurchase: c.MinimumPurchase,
		IsActive:        *c.IsActive,
		ExpiresAt:       c.ExpiresAt,
	}
}

func BuildCouponWithRemainingListResponse(coupons []models.Coupon) CouponWithRemainingListResponse {
	res := []CouponWithRemainingResponse{}

	for _, c := range coupons {
		res = append(res, BuildCouponWithRemainingResponse(c))
	}

	return CouponWithRemainingListResponse{
		Entries: res,
	}
}
