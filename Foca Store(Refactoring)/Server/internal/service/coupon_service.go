package service

import (
	"errors"
	"time"
	"voca-store/internal/domain/models"
	"voca-store/internal/domain/repository"
	"voca-store/internal/domain/dto/request"
	"voca-store/internal/domain/dto/response"
)

type CouponService interface {
	Create(req request.CreateCouponRequest) (response.CouponResponse, error)
	GetCoupons() (response.CouponWithRemainingListResponse, error)
	Update(id uint, req request.UpdateCouponRequest) (response.CouponResponse, error)
	Delete(id uint) error
	
	ClaimCoupon(userID uint, couponID uint) (response.UserCouponResponse, error)
	GetMyCoupons(userID uint) (response.UserCouponListResponse, error)
	RemoveCoupon(userID uint, userCouponID uint) error
}

type couponService struct {
	couponRepo repository.CouponRepository
}

func NewCouponService(couponRepo repository.CouponRepository) CouponService {
	return &couponService{couponRepo}
}

func (s *couponService) Create(req request.CreateCouponRequest) (response.CouponResponse, error) {
	coupon := models.Coupon{
		Code:      req.Code,
		Type:      req.Type,
		Value:     req.Value,
		Quota:     req.Quota,
		IsActive:  &req.ISActive,
		ExpiresAt: &req.ExpiresAt,
	}

	if err := s.couponRepo.Create(&coupon); err != nil {
		return response.CouponResponse{}, err
	}

	return response.BuildCouponResponse(coupon), nil
}

func (s *couponService) GetCoupons() (response.CouponWithRemainingListResponse, error) {
	coupons, err := s.couponRepo.FindAll()
	if err != nil {
		return response.CouponWithRemainingListResponse{}, err
	}
	return response.BuildCouponWithRemainingListResponse(coupons), nil
}

func (s *couponService) Update(id uint, req request.UpdateCouponRequest) (response.CouponResponse, error) {
	coupon, err := s.couponRepo.FindByID(id)
	if err != nil {
		return response.CouponResponse{}, err
	}

	updates := make(map[string]interface{})
	if req.Code != nil { updates["code"] = *req.Code }
	if req.Type != nil { updates["type"] = *req.Type }
	if req.Value != nil { updates["value"] = *req.Value }
	if req.Quota != nil { updates["quota"] = *req.Quota }
	if req.ISActive != nil { updates["is_active"] = *req.ISActive }
	if req.ExpiresAt != nil { updates["expire_at"] = *req.ExpiresAt }

	if err := s.couponRepo.Update(&coupon, updates); err != nil {
		return response.CouponResponse{}, err
	}
	
	coupon, _ = s.couponRepo.FindByID(id)
	return response.BuildCouponResponse(coupon), nil
}

func (s *couponService) Delete(id uint) error {
	coupon, err := s.couponRepo.FindByID(id)
	if err != nil { return err }
	return s.couponRepo.Delete(&coupon)
}

func (s *couponService) ClaimCoupon(userID uint, couponID uint) (response.UserCouponResponse, error) {
	var userCoupon models.UserCoupon
	err := s.couponRepo.WithTransaction(func(txRepo repository.CouponRepository) error {
		coupon, err := txRepo.FindByID(couponID)
		if err != nil { return errors.New("coupon not found") }

		if coupon.IsActive != nil && !*coupon.IsActive { return errors.New("coupon is not active") }
		if coupon.ExpiresAt != nil && coupon.ExpiresAt.Before(time.Now()) { return errors.New("coupon has expired") }
		if coupon.Quota > 0 && coupon.UsedCount >= coupon.Quota { return errors.New("coupon quota exceeded") }

		_, err = txRepo.FindUserCoupon(userID, coupon.ID)
		if err == nil { return errors.New("coupon already claimed") }

		if coupon.Quota > 0 {
			txRepo.IncrementUsedCount(coupon.ID, 1)
		}

		userCoupon = models.UserCoupon{
			UserID:   userID,
			CouponID: coupon.ID,
		}
		if err := txRepo.CreateUserCoupon(&userCoupon); err != nil { return err }
		
		// Reload for detail response
		userCoupon.Coupon = coupon
		return nil
	})

	if err != nil { return response.UserCouponResponse{}, err }
	return response.BuildUserCouponResponse(userCoupon), nil
}

func (s *couponService) GetMyCoupons(userID uint) (response.UserCouponListResponse, error) {
	userCoupons, err := s.couponRepo.FindMyCoupons(userID)
	if err != nil { return response.UserCouponListResponse{}, err }
	return response.BuildUserCouponListResponse(userCoupons), nil
}

func (s *couponService) RemoveCoupon(userID uint, userCouponID uint) error {
	return s.couponRepo.WithTransaction(func(txRepo repository.CouponRepository) error {
		userCoupon, err := txRepo.FindUserCouponByID(userCouponID, userID)
		if err != nil { return errors.New("coupon not found") }
		if userCoupon.UsedAt != nil { return errors.New("cannot remove used coupon") }

		coupon, _ := txRepo.FindByID(userCoupon.CouponID)
		if coupon.Quota > 0 && coupon.UsedCount > 0 {
			txRepo.IncrementUsedCount(coupon.ID, -1)
		}

		return txRepo.DeleteUserCoupon(&userCoupon)
	})
}
