package service

import (
	"context"
	"errors"
	"fmt"
	"voca-plane/internal/domain/dto/request"
	"voca-plane/internal/domain/dto/response"
	"voca-plane/internal/domain/models"
	"voca-plane/internal/repository"
	"voca-plane/pkg/helper"

	"gorm.io/gorm"
)

type AdminService struct {
	adminRepo   repository.AdminRepository
	userRepo    repository.UserRepository
	flightRepo  repository.FlightRepository
	airlineRepo repository.AirlineRepository
	airportRepo repository.AirportRepository
	promoRepo   repository.PromoRepository
	db          *gorm.DB
}

func NewAdminService(
	adminRepo repository.AdminRepository,
	userRepo repository.UserRepository,
	flightRepo repository.FlightRepository,
	airlineRepo repository.AirlineRepository,
	airportRepo repository.AirportRepository,
	promoRepo repository.PromoRepository,
	db *gorm.DB,
) *AdminService {
	return &AdminService{
		adminRepo:   adminRepo,
		userRepo:    userRepo,
		flightRepo:  flightRepo,
		airlineRepo: airlineRepo,
		airportRepo: airportRepo,
		promoRepo:   promoRepo,
		db:          db,
	}
}

func (s *AdminService) GetDashboardStats(ctx context.Context) (*repository.DashboardStats, error) {
	return s.adminRepo.GetDashboardStats(ctx)
}

func (s *AdminService) GetAllUsers(ctx context.Context, page, limit int, sortBy, order string) ([]response.UserResponse, int64, error) {
	users, total, err := s.adminRepo.GetAllUsers(ctx, page, limit, sortBy, order)
	if err != nil {
		return nil, 0, err
	}

	res := make([]response.UserResponse, 0)
	for _, u := range users {
		res = append(res, response.ToUserResponse(u))
	}
	return res, total, nil
}

func (s *AdminService) UpdateUserRole(ctx context.Context, userID uint, role string) error {
	validRoles := map[string]bool{
		models.RoleUser:       true,
		models.RoleAdmin:      true,
		models.RoleSuperAdmin: true,
	}
	if !validRoles[role] {
		return errors.New("invalid role")
	}

	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := s.adminRepo.UpdateUserRole(ctx, tx, userID, role); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *AdminService) DeleteUser(ctx context.Context, id uint) error {
	return s.userRepo.Delete(ctx, id)
}

func (s *AdminService) RestoreUser(ctx context.Context, id uint) error {
	return s.userRepo.Restore(ctx, id)
}

func (s *AdminService) BanUser(ctx context.Context, id uint, reason string) error {
	return s.userRepo.Ban(ctx, id, reason)
}

func (s *AdminService) UnbanUser(ctx context.Context, id uint) error {
	return s.userRepo.Unban(ctx, id)
}

func (s *AdminService) GetAllTransactions(ctx context.Context, page, limit int, sortBy, order string) ([]response.TransactionResponse, int64, error) {
	transactions, total, err := s.adminRepo.GetAllTransactions(ctx, page, limit, sortBy, order)
	if err != nil {
		return nil, 0, err
	}

	var res []response.TransactionResponse
	for _, t := range transactions {
		res = append(res, response.ToTransactionResponse(t))
	}
	return res, total, nil
}

func (s *AdminService) GetAllFlights(ctx context.Context, page, limit int, sortBy, order string) ([]response.FlightResponse, int64, error) {
	flights, total, err := s.flightRepo.GetAll(ctx, page, limit, sortBy, order)
	if err != nil {
		return nil, 0, err
	}

	flightResponses := make([]response.FlightResponse, 0)
	for _, f := range flights {
		flightResponses = append(flightResponses, response.ToFlightResponse(f))
	}

	return flightResponses, total, nil
}

func (s *AdminService) CreateFlight(ctx context.Context, flight *models.Flight, classCount int, classPrices []request.ClassPriceRequest) (*response.FlightResponse, error) {
	if err := helper.ValidateFlightInput(flight, classCount, classPrices); err != nil {
		return nil, err
	}

	classMap := helper.MapClassPrices(classPrices)

	allocations := helper.CalculateSeatAllocation(
		flight.TotalSeats,
		classMap,
		classCount,
	)

	tx := s.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Create flight
	if err := s.flightRepo.Create(ctx, tx, flight); err != nil {
		tx.Rollback()
		return nil, err
	}

	// 2. Create flight classes
	for _, alloc := range allocations {
		fClass := &models.FlightClass{
			FlightID:  flight.ID,
			ClassType: alloc.ClassType,
			Price:     alloc.Price,
		}

		if err := s.flightRepo.CreateClass(ctx, tx, fClass); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// 3. Generate seat codes and upsert master Seat rows
	seatCodes := helper.GenerateSeatCodes(flight.TotalRows, flight.TotalColumns)

	seats, err := s.flightRepo.GetOrCreateSeats(ctx, tx, seatCodes)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Sort seats to match the order of seatCodes
	seatMap := make(map[string]models.Seat)
	for _, s := range seats {
		seatMap[s.SeatCode] = s
	}
	orderedSeats := make([]models.Seat, 0, len(seatCodes))
	for _, code := range seatCodes {
		if s, ok := seatMap[code]; ok {
			orderedSeats = append(orderedSeats, s)
		}
	}

	// 4. Generate FlightSeat pivot entries distributed across allocations
	flightSeats := helper.GenerateFlightSeatModels(flight.ID, orderedSeats, allocations)

	if err := s.flightRepo.BulkCreateFlightSeats(ctx, tx, flightSeats); err != nil {
		tx.Rollback()
		return nil, err
	}

	// 5. Reload flight with all relations
	flight, err = s.flightRepo.GetFlightWithRelations(ctx, tx, flight.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	res := response.ToFlightResponse(*flight)

	return &res, nil
}

func (s *AdminService) UpdateFlight(ctx context.Context, flight *models.Flight) (*response.FlightResponse, error) {

	tx := s.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Get existing flight
	oldFlight, err := s.flightRepo.GetFlightWithClasses(ctx, tx, flight.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// 1a. Prepare classPrices dari oldFlight untuk validasi
	classPrices := []request.ClassPriceRequest{}
	for _, fc := range oldFlight.FlightClasses {
		classPrices = append(classPrices, request.ClassPriceRequest{
			ClassType: fc.ClassType,
			Price:     fc.Price,
		})
	}

	// 1b. Validasi flight input (rows, seats, classCount, classPrices)
	if err := helper.ValidateFlightInput(flight, len(oldFlight.FlightClasses), classPrices); err != nil {
		tx.Rollback()
		return nil, err
	}

	// 2. Validate seat capacity masih optional tambahan
	maxCapacity := flight.TotalRows * flight.TotalColumns
	if flight.TotalSeats > maxCapacity {
		tx.Rollback()
		return nil, fmt.Errorf(
			"total_seats (%d) exceeds rows × columns (%d)",
			flight.TotalSeats,
			maxCapacity,
		)
	}

	// 3. Update flight
	if err := s.flightRepo.Update(ctx, tx, flight); err != nil {
		tx.Rollback()
		return nil, err
	}

	// 4. Detect seat layout change
	seatsChanged :=
		oldFlight.TotalSeats != flight.TotalSeats ||
			oldFlight.TotalRows != flight.TotalRows ||
			oldFlight.TotalColumns != flight.TotalColumns

	if seatsChanged && len(oldFlight.FlightClasses) > 0 {

		classCount := len(oldFlight.FlightClasses)

		classMap := map[string]float64{}
		for _, fc := range oldFlight.FlightClasses {
			classMap[fc.ClassType] = fc.Price
		}

		allocations := helper.CalculateSeatAllocation(
			flight.TotalSeats,
			classMap,
			classCount,
		)

		// 5. Delete old flight seats
		if err := s.flightRepo.DeleteFlightSeatsByFlightID(ctx, tx, flight.ID); err != nil {
			tx.Rollback()
			return nil, err
		}

		// 6. Generate new seat codes dan upsert master Seat rows
		seatCodes := helper.GenerateSeatCodes(flight.TotalRows, flight.TotalColumns)

		seats, err := s.flightRepo.GetOrCreateSeats(ctx, tx, seatCodes)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		// Sort seats untuk match order seatCodes
		seatMap := make(map[string]models.Seat)
		for _, s := range seats {
			seatMap[s.SeatCode] = s
		}
		orderedSeats := make([]models.Seat, 0, len(seatCodes))
		for _, code := range seatCodes {
			if s, ok := seatMap[code]; ok {
				orderedSeats = append(orderedSeats, s)
			}
		}

		// 7. Regenerate FlightSeat pivot entries
		flightSeats := helper.GenerateFlightSeatModels(flight.ID, orderedSeats, allocations)

		if err := s.flightRepo.BulkCreateFlightSeats(ctx, tx, flightSeats); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// 8. Reload flight
	flight, err = s.flightRepo.GetFlightWithRelations(ctx, tx, flight.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	res := response.ToFlightResponse(*flight)

	return &res, nil
}

func (s *AdminService) DeleteFlight(ctx context.Context, id uint) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := s.flightRepo.Delete(ctx, tx, id); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *AdminService) GetAllAirlines(ctx context.Context, page, limit int, sortBy, order string) ([]response.AirlineResponse, int64, error) {
	airlines, total, err := s.airlineRepo.GetAll(ctx, page, limit, sortBy, order)
	if err != nil {
		return nil, 0, err
	}

	res := make([]response.AirlineResponse, 0)
	for _, a := range airlines {
		res = append(res, response.ToAirlineResponse(a))
	}
	return res, total, nil
}

func (s *AdminService) CreateAirline(ctx context.Context, airline *models.Airline) (*response.AirlineResponse, error) {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := s.airlineRepo.Create(ctx, tx, airline); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	res := response.ToAirlineResponse(*airline)
	return &res, nil
}

func (s *AdminService) UpdateAirline(ctx context.Context, airline *models.Airline) (*response.AirlineResponse, error) {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := s.airlineRepo.Update(ctx, tx, airline); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	res := response.ToAirlineResponse(*airline)
	return &res, nil
}

func (s *AdminService) DeleteAirline(ctx context.Context, id uint) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	airline, err := s.airlineRepo.GetByID(ctx, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	if airline.LogoPublicID != "" {
		helper.DeleteImage(airline.LogoPublicID)
	}

	if err := s.airlineRepo.Delete(ctx, tx, id); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *AdminService) GetAllAirports(ctx context.Context, page, limit int, sortBy, order string) ([]response.AirportResponse, int64, error) {
	airports, total, err := s.airportRepo.GetAll(ctx, page, limit, sortBy, order)
	if err != nil {
		return nil, 0, err
	}

	res := make([]response.AirportResponse, 0)
	for _, a := range airports {
		res = append(res, response.ToAirportResponse(a))
	}
	return res, total, nil
}

func (s *AdminService) CreateAirport(ctx context.Context, airport *models.Airport) (*response.AirportResponse, error) {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := s.airportRepo.Create(ctx, tx, airport); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	res := response.ToAirportResponse(*airport)
	return &res, nil
}

func (s *AdminService) UpdateAirport(ctx context.Context, airport *models.Airport) (*response.AirportResponse, error) {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := s.airportRepo.Update(ctx, tx, airport); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	res := response.ToAirportResponse(*airport)
	return &res, nil
}

func (s *AdminService) DeleteAirport(ctx context.Context, id uint) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := s.airportRepo.Delete(ctx, tx, id); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *AdminService) GetAllPromos(ctx context.Context, page, limit int, sortBy, order string) ([]response.PromoResponse, int64, error) {
	promos, total, err := s.promoRepo.GetAll(ctx, page, limit, sortBy, order)
	if err != nil {
		return nil, 0, err
	}

	var res []response.PromoResponse
	for _, p := range promos {
		res = append(res, response.ToPromoResponse(p))
	}
	return res, total, nil
}

func (s *AdminService) CreatePromo(ctx context.Context, promo *models.PromoCode) (*response.PromoResponse, error) {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := s.promoRepo.Create(ctx, tx, promo); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	res := response.ToPromoResponse(*promo)
	return &res, nil
}

func (s *AdminService) UpdatePromo(ctx context.Context, promo *models.PromoCode) (*response.PromoResponse, error) {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := s.promoRepo.Update(ctx, tx, promo); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	res := response.ToPromoResponse(*promo)
	return &res, nil
}

func (s *AdminService) DeletePromo(ctx context.Context, id uint) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := s.promoRepo.Delete(ctx, tx, id); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *AdminService) GetFlightByID(ctx context.Context, id uint) (*response.FlightResponse, error) {
	flight, err := s.flightRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	response := response.ToFlightResponse(*flight)
	return &response, nil
}

func (s *AdminService) GetFlightModelByID(ctx context.Context, id uint) (*models.Flight, error) {
	return s.flightRepo.GetByID(ctx, id)
}

func (s *AdminService) GetAirlineByID(ctx context.Context, id uint) (*models.Airline, error) {
	return s.airlineRepo.GetByID(ctx, id)
}

func (s *AdminService) GetAirportByID(ctx context.Context, id uint) (*models.Airport, error) {
	return s.airportRepo.GetByID(ctx, id)
}

func (s *AdminService) GetPromoByID(ctx context.Context, id uint) (*models.PromoCode, error) {
	return s.promoRepo.GetByID(ctx, id)
}
