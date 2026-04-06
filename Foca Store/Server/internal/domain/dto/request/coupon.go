package request

import "time"

type CreateCouponRequest struct {
	Code      string    `json:"code" binding:"required"`
	Type      string    `json:"type" binding:"required,oneof=percentage fixed"`
	Value     float64   `json:"value" binding:"required,gt=0"`
	Quota     int       `json:"quota" binding:"required,gt=0"`
	ISActive  bool      `json:"is_active" binding:"required"`
	ExpiresAt time.Time `json:"expires_at" binding:"required"`
}

type UpdateCouponRequest struct {
	Code      *string    `json:"code"`
	Type      *string    `json:"type" binding:"omitempty,oneof=percentage fixed"`
	Value     *float64   `json:"value" binding:"omitempty,gt=0"`
	Quota     *int       `json:"quota" binding:"omitempty,gt=0"`
	ISActive  *bool      `json:"is_active" binding:"omitempty"`
	ExpiresAt *time.Time `json:"expires_at" binding:"omitempty"`
}
