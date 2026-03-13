package service

import (
	"context"
	"errors"
	"fmt"
	"time"
	"voca-plane/internal/domain/dto/request"
	"voca-plane/internal/domain/dto/response"
	"voca-plane/internal/domain/models"
	"voca-plane/internal/repository"
	"voca-plane/pkg/helper"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionService struct {
	txRepo     repository.TransactionRepository
	flightRepo repository.FlightRepository
	promoRepo  repository.PromoRepository
	db         *gorm.DB
	midtrans   *helper.MidTransClient
}

func NewTransactionService(txRepo repository.TransactionRepository, flightRepo repository.FlightRepository, promoRepo repository.PromoRepository, db *gorm.DB, midtrans *helper.MidTransClient) *TransactionService {
	return &TransactionService{
		txRepo:     txRepo,
		flightRepo: flightRepo,
		promoRepo:  promoRepo,
		db:         db,
		midtrans:   midtrans,
	}
}

func (s *TransactionService) CreateTransaction(ctx context.Context, userID uint, req request.CreateTransactionRequest) (*response.TransactionResponse, error) {
	lockDuration := 10 * time.Minute

	// Validate no duplicate seat codes
	seatCodeMap := map[string]bool{}
	seatCodes := make([]string, 0, len(req.Passengers))
	for _, p := range req.Passengers {
		if seatCodeMap[p.SeatNumber] {
			return nil, fmt.Errorf("duplicated seat_number %s", p.SeatNumber)
		}
		seatCodeMap[p.SeatNumber] = true
		seatCodes = append(seatCodes, p.SeatNumber)
	}

	var transaction models.Transaction
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Fetch seats with FOR UPDATE lock by seat codes
		lockedSeats, err := s.flightRepo.GetFlightSeatsByCodes(ctx, tx, req.FlightID, seatCodes)
		if err != nil {
			return fmt.Errorf("failed to fetch seats: %w", err)
		}

		if len(lockedSeats) != len(seatCodes) {
			return errors.New("one or more seats not found in this flight")
		}

		// Pre-fetch all classes for this flight to get prices and IDs
		flight, err := s.flightRepo.GetFlightWithClasses(ctx, tx, req.FlightID)
		if err != nil {
			return fmt.Errorf("failed to fetch flight classes: %w", err)
		}

		classMap := make(map[string]models.FlightClass)
		for _, c := range flight.FlightClasses {
			classMap[c.ClassType] = c
		}

		seatIDs := make([]uint, 0, len(lockedSeats))
		items := make([]models.TransactionItem, len(req.Passengers))
		subtotal := 0.0

		for i, p := range req.Passengers {
			var currentSeat models.FlightSeat
			found := false
			for _, ls := range lockedSeats {
				if ls.Seat.SeatCode == p.SeatNumber {
					currentSeat = ls
					found = true
					break
				}
			}

			if !found || !currentSeat.IsAvailable {
				return fmt.Errorf("seat %s is not available", p.SeatNumber)
			}

			class, ok := classMap[currentSeat.ClassType]
			if !ok {
				return fmt.Errorf("class %s not found for seat %s", currentSeat.ClassType, p.SeatNumber)
			}

			seatIDs = append(seatIDs, currentSeat.ID)
			subtotal += class.Price

			items[i] = models.TransactionItem{
				PassengerName: p.FullName,
				Nationality:   p.Nationality,
				PassportNo:    p.PassportNo,
				SeatNumber:    p.SeatNumber,
				FlightSeatID:  currentSeat.ID,
				FlightClassID: class.ID,
				Price:         class.Price,
			}
		}

		discount := 0.0
		if req.PromoCode != nil {
			promo, err := s.promoRepo.GetByCode(ctx, *req.PromoCode)
			if err == nil && promo.IsActive {
				discount = subtotal * (promo.Discount / 100)
			}
		}

		totalPrice := subtotal - discount

		code := uuid.New().String()
		transaction = models.Transaction{
			Code:          code,
			UserID:        userID,
			FlightID:      req.FlightID,
			TotalPrice:    totalPrice,
			PaymentStatus: "PENDING",
			PromoCode:     req.PromoCode,
			Discount:      discount,
			ExpiresAt:     time.Now().Add(lockDuration),
		}

		if err := s.txRepo.Create(ctx, tx, &transaction); err != nil {
			return err
		}

		// Lock seats with transaction reference
		if err := s.flightRepo.LockSeats(ctx, tx, seatIDs, transaction.ID, time.Now().Add(lockDuration)); err != nil {
			return err
		}

		// Set transaction ID for items
		for i := range items {
			items[i].TransactionID = transaction.ID
		}

		if err := s.txRepo.CreateTransactionItems(ctx, tx, items); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if transaction.TotalPrice <= 0 {
		transaction.PaymentStatus = "PAID"
		transactionRes, err := s.txRepo.GetByCode(ctx, transaction.Code)
		if err != nil {
			return nil, err
		}
		resDto := response.ToTransactionResponse(*transactionRes)
		return &resDto, nil
	}

	res, err := s.midtrans.CreatePayment(
		transaction.Code,
		transaction.TotalPrice,
	)

	if err != nil {
		return nil, err
	}

	transaction.PaymentURL = res.RedirectURL

	if err := s.txRepo.UpdatePaymentURL(ctx, transaction.Code, res.RedirectURL); err != nil {
		return nil, err
	}

	transactionRes, err := s.txRepo.GetByCode(ctx, transaction.Code)
	if err != nil {
		return nil, err
	}

	resDto := response.ToTransactionResponse(*transactionRes)
	return &resDto, nil
}

func (s *TransactionService) PayTransaction(ctx context.Context, code string) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	transaction, err := s.txRepo.GetByCode(ctx, code)
	if err != nil {
		tx.Rollback()
		return errors.New("transaction not found")
	}

	if transaction.PaymentStatus == "PAID" {
		tx.Rollback()
		return errors.New("already paid")
	}

	if err := s.txRepo.UpdatePaymentStatus(ctx, tx, transaction.ID, "PAID"); err != nil {
		tx.Rollback()
		return err
	}

	// Finalize seats: clear lock timer but keep them booked
	if err := s.flightRepo.FinalizeSeats(ctx, tx, transaction.ID); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *TransactionService) GetUserTransactions(ctx context.Context, userID uint, page, limit int) ([]response.TransactionResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	transactions, total, err := s.txRepo.GetByUserID(ctx, userID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	var res []response.TransactionResponse
	for _, t := range transactions {
		res = append(res, response.ToTransactionResponse(t))
	}

	return res, total, nil
}

func (s *TransactionService) GetUserTransactionsAll(ctx context.Context, userID uint) ([]response.TransactionResponse, error) {
	transactions, err := s.txRepo.GetByUserIDAll(ctx, userID)
	if err != nil {
		return nil, err
	}

	var res []response.TransactionResponse
	for _, t := range transactions {
		res = append(res, response.ToTransactionResponse(t))
	}

	return res, nil
}

func (s *TransactionService) GetTransactionByCode(ctx context.Context, code string) (*response.TransactionResponse, error) {
	transaction, err := s.txRepo.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	if transaction.PaymentStatus == "PENDING" && time.Now().After(transaction.ExpiresAt) {
		err := s.ExpireTransaction(ctx, code)
		if err != nil {
			return nil ,err
		}

		transaction, err = s.txRepo.GetByCode(ctx, code)
		if err != nil {
			return nil, err
		}
	}

	res := response.ToTransactionResponse(*transaction)
	return &res, nil
}

func (s *TransactionService) CancelTransaction(ctx context.Context, userID uint, code string) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	transaction, err := s.txRepo.GetByCode(ctx, code)
	if err != nil {
		tx.Rollback()
		return errors.New("transaction not found")
	}

	if transaction.UserID != userID {
		tx.Rollback()
		return errors.New("unauthorized")
	}

	if transaction.PaymentStatus == "PAID" {
		tx.Rollback()
		return errors.New("cannot cancel paid transaction")
	}

	// Release locked seats back to available
	if err := s.flightRepo.ReleaseSeats(ctx, tx, transaction.ID); err != nil {
		tx.Rollback()
		return err
	}

	if err := s.txRepo.Delete(ctx, tx, code); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *TransactionService) ExpireTransaction(ctx context.Context, code string) error {
	tx := s.db.Begin()

	transaction, err := s.txRepo.GetByCode(ctx, code)
	if err != nil {
		tx.Rollback()
		return err
	}

	if transaction.PaymentStatus != "PENDING" {
		tx.Rollback()
		return nil
	}

	// if time.Now().Before(transaction.ExpiresAt) {
	// 	tx.Rollback()
	// 	return nil
	// }

	if err := s.txRepo.UpdatePaymentStatus(ctx, tx, transaction.ID, "EXPIRED"); err != nil {
		tx.Rollback()
		return err
	}

	if err := s.flightRepo.ReleaseSeats(ctx, tx, transaction.ID); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}