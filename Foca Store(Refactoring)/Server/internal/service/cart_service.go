package service

import (
	"errors"
	"voca-store/internal/domain/models"
	"voca-store/internal/domain/repository"
	"voca-store/internal/domain/dto/request"
	"voca-store/internal/domain/dto/response"

	"gorm.io/gorm"
)

type CartService interface {
	GetCart(userID uint) (response.CartResponse, error)
	AddToCart(userID uint, req request.AddToCartRequest) error
	RemoveItem(userID uint, itemID uint) error
	RemoveItems(userID uint, req request.RemoveCartItemRequest) error
	ClearCart(userID uint) error
}

type cartService struct {
	cartRepo repository.CartRepository
	productRepo repository.ProductRepository
}

func NewCartService(cartRepo repository.CartRepository, productRepo repository.ProductRepository) CartService {
	return &cartService{cartRepo, productRepo}
}

func (s *cartService) GetCart(userID uint) (response.CartResponse, error) {
	cart, err := s.cartRepo.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create cart if not exists
			cart = models.Cart{UserID: userID}
			s.cartRepo.Create(&cart)
			cart, _ = s.cartRepo.FindByUserID(userID)
		} else {
			return response.CartResponse{}, err
		}
	}

	var cartItemResponses []response.CartItemResponse
	for _, item := range cart.Items {
		if item.Product == nil { continue }
		productResp := response.BuildProductResponse(*item.Product)
		cartItemResp := response.BuildCartItemResponse(item.ID, item.ProductID, productResp, item.Quantity, item.CreatedAt, item.UpdatedAt)
		cartItemResponses = append(cartItemResponses, cartItemResp)
	}

	return response.BuildCartResponse(cart.ID, cart.UserID, cartItemResponses, cart.CreatedAt, cart.UpdatedAt), nil
}

func (s *cartService) AddToCart(userID uint, req request.AddToCartRequest) error {
	product, err := s.productRepo.FindByID(req.ProductID)
	if err != nil {
		return errors.New("product not found")
	}

	if product.Stock < req.Quantity {
		return errors.New("insufficient stock")
	}

	cart, err := s.cartRepo.FindByUserID(userID)
	if err != nil {
		cart = models.Cart{UserID: userID}
		s.cartRepo.Create(&cart)
	}

	// Check if already in cart
	existingItem, err := s.cartRepo.FindItem(cart.ID, req.ProductID)
	if err == nil {
		newQuantity := existingItem.Quantity + req.Quantity
		if newQuantity > product.Stock {
			return errors.New("insufficient stock")
		}
		return s.cartRepo.UpdateItem(&existingItem, map[string]interface{}{"quantity": newQuantity})
	}

	// New item
	item := models.CartItem{
		CartID:    cart.ID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	}
	return s.cartRepo.AddItem(&item)
}

func (s *cartService) RemoveItem(userID uint, itemID uint) error {
	cart, _ := s.cartRepo.FindByUserID(userID)
	item, err := s.cartRepo.FindItemByID(itemID)
	if err != nil {
		return errors.New("item not found")
	}

	if item.CartID != cart.ID {
		return errors.New("access denied")
	}

	return s.cartRepo.DeleteItem(&item)
}

func (s *cartService) RemoveItems(userID uint, req request.RemoveCartItemRequest) error {
	cart, _ := s.cartRepo.FindByUserID(userID)
	
	count, _ := s.cartRepo.CountValidItems(cart.ID, req.CartItemIDs)
	if count != int64(len(req.CartItemIDs)) {
		return errors.New("some items are invalid")
	}

	return s.cartRepo.DeleteItems(cart.ID, req.CartItemIDs)
}

func (s *cartService) ClearCart(userID uint) error {
	cart, _ := s.cartRepo.FindByUserID(userID)
	return s.cartRepo.ClearCart(cart.ID)
}
