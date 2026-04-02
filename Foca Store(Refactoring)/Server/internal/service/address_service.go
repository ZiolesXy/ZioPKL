package service

import (
	"errors"
	"voca-store/internal/domain/models"
	"voca-store/internal/domain/repository"
	"voca-store/internal/domain/dto/request"
	"voca-store/internal/domain/dto/response"
	"voca-store/internal/helper"
)

type AddressService interface {
	Create(userID uint, req request.CreateAddressRequest) (response.AddressResponse, error)
	GetMyAddresses(userID uint) (response.AddressListResponse, error)
	GetByUID(uid string, userID uint) (response.AddressResponse, error)
	Update(uid string, userID uint, req request.UpdateAddressRequest) (response.AddressResponse, error)
	Delete(uid string, userID uint) error
}

type addressService struct {
	addressRepo repository.AddressRepository
}

func NewAddressService(addressRepo repository.AddressRepository) AddressService {
	return &addressService{addressRepo}
}

func (s *addressService) Create(userID uint, req request.CreateAddressRequest) (response.AddressResponse, error) {
	var finalAddress models.Address
	err := s.addressRepo.WithTransaction(func(txRepo repository.AddressRepository) error {
		uid, _ := helper.NewGenerateAddressUID(nil) // Fixed if helper doesn't need DB

		address := models.Address{
			UID:           uid,
			UserID:        userID,
			Label:         req.Label,
			RecipientName: req.RecipientName,
			Phone:         req.Phone,
			AddressLine:   req.AddressLine,
			City:          req.City,
			Province:      req.Province,
			PostalCode:    req.PostalCode,
			IsPrimary:     req.IsPrimary,
		}

		if address.IsPrimary {
			txRepo.UnsetPrimary(userID)
		}

		if err := txRepo.Create(&address); err != nil {
			return err
		}
		finalAddress = address
		return nil
	})

	if err != nil {
		return response.AddressResponse{}, err
	}
	return response.BuildAddressResponse(finalAddress), nil
}

func (s *addressService) GetMyAddresses(userID uint) (response.AddressListResponse, error) {
	addresses, err := s.addressRepo.FindAll(userID)
	if err != nil {
		return response.AddressListResponse{}, err
	}
	return response.BuildAddressListResponse(addresses), nil
}

func (s *addressService) GetByUID(uid string, userID uint) (response.AddressResponse, error) {
	address, err := s.addressRepo.FindByUID(uid, userID)
	if err != nil {
		return response.AddressResponse{}, errors.New("address not found")
	}
	return response.BuildAddressResponse(address), nil
}

func (s *addressService) Update(uid string, userID uint, req request.UpdateAddressRequest) (response.AddressResponse, error) {
	var finalAddress models.Address
	err := s.addressRepo.WithTransaction(func(txRepo repository.AddressRepository) error {
		address, err := txRepo.FindByUID(uid, userID)
		if err != nil {
			return errors.New("address not found")
		}

		updates := make(map[string]interface{})
		if req.Label != nil { updates["label"] = *req.Label }
		if req.RecipientName != nil { updates["recipient_name"] = *req.RecipientName }
		if req.Phone != nil { updates["phone"] = *req.Phone }
		if req.AddressLine != nil { updates["address_line"] = *req.AddressLine }
		if req.City != nil { updates["city"] = *req.City }
		if req.Province != nil { updates["province"] = *req.Province }
		if req.PostalCode != nil { updates["postal_code"] = *req.PostalCode }

		if req.IsPrimary != nil {
			if *req.IsPrimary {
				txRepo.UnsetPrimary(userID)
			}
			updates["is_primary"] = *req.IsPrimary
		}

		if err := txRepo.Update(&address, updates); err != nil {
			return err
		}
		finalAddress, _ = txRepo.FindByID(address.ID, userID)
		return nil
	})

	if err != nil {
		return response.AddressResponse{}, err
	}
	return response.BuildAddressResponse(finalAddress), nil
}

func (s *addressService) Delete(uid string, userID uint) error {
	address, err := s.addressRepo.FindByUID(uid, userID)
	if err != nil {
		return errors.New("address not found")
	}
	return s.addressRepo.Delete(&address)
}
