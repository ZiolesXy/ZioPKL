package service

import (
	"errors"
	"fmt"
	"time"
	"voca-store/internal/domain/models"
	"voca-store/internal/domain/repository"
	"voca-store/internal/domain/dto/request"
	"voca-store/internal/domain/dto/response"
	"voca-store/internal/helper"
)

type CheckoutService interface {
	Checkout(userID uint, req request.CheckoutRequest) (response.CheckoutDetailResponse, error)
	ApproveCheckout(id uint) (response.CheckoutDetailResponse, error)
	RejectCheckout(id uint) (response.CheckoutDetailResponse, error)
	GetCheckouts() (response.CheckoutListResponse, error)
	GetMyCheckouts(userID uint) (response.CheckoutListResponse, error)
	GetCheckoutByUID(uid string, userID uint) (response.CheckoutDetailResponse, error)
	DeleteMyCheckout(uid string, userID uint) error
	HandleMidtransWebhook(payload map[string]interface{}) error
}

type checkoutService struct {
	checkoutRepo repository.CheckoutRepository
	userRepo     repository.UserRepository
}

func NewCheckoutService(checkoutRepo repository.CheckoutRepository, userRepo repository.UserRepository) CheckoutService {
	return &checkoutService{checkoutRepo, userRepo}
}

func (s *checkoutService) Checkout(userID uint, req request.CheckoutRequest) (response.CheckoutDetailResponse, error) {
	var finalCheckout models.Checkout

	err := s.checkoutRepo.WithTransaction(func(txRepo repository.CheckoutRepository) error {
		cart, err := txRepo.FindCartByUserID(userID)
		if err != nil {
			return errors.New("cart not found")
		}

		items, err := txRepo.FindCartItems(cart.ID, req.CartItemIDs)
		if err != nil || len(items) != len(req.CartItemIDs) {
			return errors.New("invalid cart items")
		}

		var subtotalCents int64 = 0
		for _, item := range items {
			if item.Product.Stock < item.Quantity {
				return fmt.Errorf("insufficient stock for %s", item.Product.Name)
			}
			txRepo.UpdateProductStock(item.ProductID, item.Quantity)
			subtotalCents += int64(item.Product.Price) * int64(item.Quantity)
		}

		discountAmountCents := int64(0)
		var coupon models.Coupon
		if req.CouponCode != nil && *req.CouponCode != "" {
			coupon, err = txRepo.FindCouponByCode(*req.CouponCode)
			if err != nil {
				return errors.New("invalid coupon")
			}
			if coupon.ExpiresAt != nil && coupon.ExpiresAt.Before(time.Now()) {
				return errors.New("coupon has expired")
			}
			if subtotalCents < int64(coupon.MinimumPurchase) {
				return errors.New("minimum purchase not reached")
			}
			userCoupon, err := txRepo.FindUserCoupon(userID, coupon.ID)
			if err != nil {
				return errors.New("coupon not claimed")
			}
			if userCoupon.UsedAt != nil {
				return errors.New("coupon already used")
			}

			if coupon.Type == "percentage" {
				discountAmountCents = subtotalCents * int64(coupon.Value) / 100
			} else {
				discountAmountCents = int64(coupon.Value)
			}

			txRepo.UpdateUserCoupon(&userCoupon, map[string]interface{}{"used_at": time.Now()})
		}

		totalCents := subtotalCents - discountAmountCents
		if totalCents < 0 { totalCents = 0 }

		address, err := txRepo.FindAddressByUID(req.AddressUID, userID)
		if err != nil {
			return errors.New("address not found")
		}

		// Generate UID (assuming helper works with tx or just returns a string)
		uid, _ := txRepo.GenerateUID()

		checkout := models.Checkout{
			UID:            uid,
			UserID:         userID,
			AddressID:      &address.ID,
			Subtotal:       float64(subtotalCents),
			DiscountAmount: float64(discountAmountCents),
			TotalPrice:     float64(totalCents),
			Status:         "pending",
		}
		if coupon.ID != 0 { checkout.CouponID = &coupon.ID }

		user, _ := s.userRepo.FindByID(userID)
		checkout.User = &user
		checkout.Address = &address

		for _, item := range items {
			checkout.Items = append(checkout.Items, models.CheckoutItem{
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Price:     item.Product.Price,
				Product:   *item.Product,
			})
		}

		checkout.WhatsappURL = helper.GenerateCheckoutWhatsappURL(checkout)

		if totalCents == 0 {
			checkout.PaymentStatus = "paid"
			checkout.Status = "approved"
			checkout.MidtransOrderID = "FREE-" + checkout.UID
		} else {
			midtransResp, err := helper.CreateSnapTransaction(
				checkout.UID, totalCents, discountAmountCents,
				user.Name, user.Email, user.TelephoneNumber, address, items,
			)
			if err != nil { return err }
			checkout.MidtransOrderID = midtransResp.OrderID
			checkout.SnapToken = midtransResp.Token
			checkout.PaymentURL = midtransResp.RedirectURL
			checkout.PaymentStatus = "pending"
		}

		if err := txRepo.Create(&checkout); err != nil { return err }

		txRepo.DeleteCartItems(cart.ID, req.CartItemIDs)
		finalCheckout = checkout
		return nil
	})

	if err != nil { return response.CheckoutDetailResponse{}, err }
	
	// Reload for full detail
	resCheckout, _ := s.checkoutRepo.FindByID(finalCheckout.ID)
	return response.BuildCheckoutDetailResponse(resCheckout), nil
}

func (s *checkoutService) ApproveCheckout(id uint) (response.CheckoutDetailResponse, error) {
	checkout, err := s.checkoutRepo.FindByID(id)
	if err != nil { return response.CheckoutDetailResponse{}, err }
	if checkout.Status != "pending" { return response.CheckoutDetailResponse{}, errors.New("invalid status") }
	
	checkout.Status = "approved"
	s.checkoutRepo.Save(&checkout)
	return response.BuildCheckoutDetailResponse(checkout), nil
}

func (s *checkoutService) RejectCheckout(id uint) (response.CheckoutDetailResponse, error) {
	var finalCheckout models.Checkout
	err := s.checkoutRepo.WithTransaction(func(txRepo repository.CheckoutRepository) error {
		checkout, err := txRepo.FindByID(id)
		if err != nil { return err }
		if checkout.Status != "pending" { return errors.New("invalid status") }
		
		for _, item := range checkout.Items {
			txRepo.UpdateProductStock(item.ProductID, -item.Quantity)
		}
		checkout.Status = "rejected"
		txRepo.Save(&checkout)
		finalCheckout = checkout
		return nil
	})
	if err != nil { return response.CheckoutDetailResponse{}, err }
	return response.BuildCheckoutDetailResponse(finalCheckout), nil
}

func (s *checkoutService) GetCheckouts() (response.CheckoutListResponse, error) {
	checkouts, err := s.checkoutRepo.FindAll()
	if err != nil {
		return response.CheckoutListResponse{}, err
	}
	return response.BuildCheckOutListResponse(checkouts), nil
}

func (s *checkoutService) GetMyCheckouts(userID uint) (response.CheckoutListResponse, error) {
	checkouts, err := s.checkoutRepo.FindByUserID(userID)
	if err != nil { return response.CheckoutListResponse{}, err }
	return response.BuildCheckOutListResponse(checkouts), nil
}

func (s *checkoutService) GetCheckoutByUID(uid string, userID uint) (response.CheckoutDetailResponse, error) {
	checkout, err := s.checkoutRepo.FindByUID(uid)
	if err != nil || checkout.UserID != userID { return response.CheckoutDetailResponse{}, errors.New("checkout not found") }
	return response.BuildCheckoutDetailResponse(checkout), nil
}

func (s *checkoutService) DeleteMyCheckout(uid string, userID uint) error {
	checkout, err := s.checkoutRepo.FindByUID(uid)
	if err != nil || checkout.UserID != userID { return errors.New("checkout not found") }
	return s.checkoutRepo.Delete(&checkout)
}

func (s *checkoutService) HandleMidtransWebhook(payload map[string]interface{}) error {
	orderID, _ := payload["order_id"].(string)
	status, _ := payload["transaction_status"].(string)

	checkout, err := s.checkoutRepo.FindByMidtransOrderID(orderID)
	if err != nil { return err }

	updates := make(map[string]interface{})
	switch status {
	case "capture", "settlement":
		updates["payment_status"] = "paid"
		updates["status"] = "approved"
	case "deny", "cancel", "expire":
		updates["payment_status"] = "failed"
		updates["status"] = "rejected"
		// Restore stock logic
		for _, item := range checkout.Items {
			s.checkoutRepo.UpdateProductStock(item.ProductID, -item.Quantity)
		}
	case "pending":
		updates["payment_status"] = "pending"
	}

	return s.checkoutRepo.Update(&checkout, updates)
}
