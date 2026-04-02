package coupon

import (
	"net/http"
	"strconv"
	"voca-store/internal/domain/dto/request"
	"voca-store/internal/domain/dto/response"
	"voca-store/internal/service"

	"github.com/gin-gonic/gin"
)

type CouponHandler struct {
	couponService service.CouponService
}

func NewCouponHandler(couponService service.CouponService) *CouponHandler {
	return &CouponHandler{couponService}
}

func (h *CouponHandler) CreateCoupon(c *gin.Context) {
	var req request.CreateCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	couponRes, err := h.couponService.Create(req)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessResponse(c, "coupon created", couponRes)
}

func (h *CouponHandler) GetCoupons(c *gin.Context) {
	coupons, err := h.couponService.GetCoupons()
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "failed to get coupons")
		return
	}

	response.SuccessResponse(c, "coupons retrieved successfully", coupons)
}

func (h *CouponHandler) UpdateCoupon(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	var req request.UpdateCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	couponRes, err := h.couponService.Update(uint(id), req)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessResponse(c, "coupon successfully updated", couponRes)
}

func (h *CouponHandler) DeleteCoupon(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	if err := h.couponService.Delete(uint(id)); err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "failed to delete")
		return
	}

	response.SuccessResponse(c, "deleted", nil)
}

func (h *CouponHandler) ClaimCoupon(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uint)
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	userCouponRes, err := h.couponService.ClaimCoupon(userID, uint(id))
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessResponse(c, "coupon claimed", userCouponRes)
}

func (h *CouponHandler) GetMyCoupons(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uint)

	myCoupons, err := h.couponService.GetMyCoupons(userID)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "failed to fetch coupons")
		return
	}

	response.SuccessResponse(c, "your coupons fetched", myCoupons)
}

func (h *CouponHandler) RemoveCoupon(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uint)
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	if err := h.couponService.RemoveCoupon(userID, uint(id)); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessResponse(c, "coupon removed", nil)
}
