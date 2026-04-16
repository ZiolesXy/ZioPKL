package service

import (
	"context"
	"voca-plane/internal/domain/dto/response"
	"voca-plane/internal/domain/models"
	"voca-plane/internal/repository"
)

type FlightService struct {
	flightRepo repository.FlightRepository
}

func NewFlightService(flightRepo repository.FlightRepository) *FlightService {
	return &FlightService{flightRepo: flightRepo}
}

func (s *FlightService) SearchFlight(ctx context.Context, origin, destination, date, classType string, page, limit int) ([]response.FlightResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	flights, total, err := s.flightRepo.Search(ctx, origin, destination, date, classType, page, limit)
	if err != nil {
		return nil, 0, err
	}

	flightResponses := make([]response.FlightResponse, 0)
	for _, f := range flights {
		flightResponses = append(flightResponses, response.ToFlightResponse(f))
	}

	return flightResponses, total, nil
}

func (s *FlightService) GetFlightByID(ctx context.Context, id uint) (*response.FlightResponse, error) {
	flight, err := s.flightRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	res := response.ToFlightResponse(*flight)
	return &res, nil
}

func (s *FlightService) GetAllFlights(ctx context.Context, page, limit int, sortBy, order string) ([]response.FlightResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
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

func (s *FlightService) GetAllFlightsFull (ctx context.Context) ([]response.FlightResponse, error) {
	flights, err := s.flightRepo.GetAllFull(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]response.FlightResponse, 0)
	for _, f := range flights {
		res = append(res, response.ToFlightResponse(f))
	}

	return res, nil
}

func (s *FlightService) CreateFlight(ctx context.Context, flight *models.Flight) error {
	return s.flightRepo.Create(ctx, nil, flight)
}

func (s *FlightService) UpdateFlight(ctx context.Context, flight *models.Flight) error {
	return s.flightRepo.Update(ctx, nil, flight)
}

func (s *FlightService) DeleteFlight(ctx context.Context, id uint) error {
	return s.flightRepo.Delete(ctx, nil, id)
}

// GetAvailableSeats returns available seats for a flight, optionally filtered by class type.
func (s *FlightService) GetAvailableSeats(ctx context.Context, flightID uint, classType string) ([]response.FlightSeatResponse, error) {
	seats, err := s.flightRepo.GetAvailableSeats(ctx, flightID, classType)
	if err != nil {
		return nil, err
	}

	var res []response.FlightSeatResponse
	for _, s := range seats {
		res = append(res, response.ToFlightSeatResponse(s))
	}

	return res, nil
}