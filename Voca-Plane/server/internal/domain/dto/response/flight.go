package response

import (
	"time"
	"voca-plane/internal/domain/models"
)

type FlightSeatResponse struct {
	ID          uint       `json:"id"`
	SeatCode    string     `json:"seat_code"`
	ClassType   string     `json:"class_type"`
	IsAvailable bool       `json:"is_available"`
	LockedUntil *time.Time `json:"locked_until,omitempty"`
}

type FlightClassResponse struct {
	ID         uint    `json:"id"`
	ClassType  string  `json:"class_type"`
	Price      float64 `json:"price"`
	TotalSeats int     `json:"total_seats,omitempty"`
}

type FlightResponse struct {
	ID             uint                  `json:"id"`
	Airline        AirlineResponse       `json:"airline"`
	Origin         AirportResponse       `json:"origin"`
	Destination    AirportResponse       `json:"destination"`
	DepartureTime  time.Time             `json:"departure_time"`
	ArrivalTime    time.Time             `json:"arrival_time"`
	FlightNumber   string                `json:"flight_number"`
	TotalSeats     int                   `json:"total_seats"`
	AvailableSeats int                   `json:"available_seats"`
	TotalRows      int                   `json:"total_rows"`
	TotalColumns   int                   `json:"total_columns"`
	FlightClasses  []FlightClassResponse `json:"classes,omitempty"`
	FlightSeats    []FlightSeatResponse  `json:"flight_seats,omitempty"`
}

func ToFlightSeatResponse(s models.FlightSeat) FlightSeatResponse {
	return FlightSeatResponse{
		ID:          s.ID,
		SeatCode:    s.Seat.SeatCode,
		ClassType:   s.ClassType,
		IsAvailable: s.IsAvailable,
		LockedUntil: s.LockedUntil,
	}
}

func ToFlightClassResponse(fc models.FlightClass, seatCount int) FlightClassResponse {
	return FlightClassResponse{
		ID:         fc.ID,
		ClassType:  fc.ClassType,
		Price:      fc.Price,
		TotalSeats: seatCount,
	}
}

func ToFlightResponse(f models.Flight) FlightResponse {
	available := f.AvailableSeats
	if available == 0 && len(f.FlightSeats) > 0 {
		for _, seat := range f.FlightSeats {
			if seat.IsAvailable {
				available++
			}
		}
	}

	// Count seats per class for the class response
	classSeatsCount := make(map[uint]int)
	for _, seat := range f.FlightSeats {
		for _, fc := range f.FlightClasses {
			if seat.ClassType == fc.ClassType {
				classSeatsCount[fc.ID]++
				break
			}
		}
	}

	var classes []FlightClassResponse
	for _, c := range f.FlightClasses {
		classes = append(classes, ToFlightClassResponse(c, classSeatsCount[c.ID]))
	}

	var seats []FlightSeatResponse
	if len(f.FlightSeats) > 0 {
		for _, s := range f.FlightSeats {
			seats = append(seats, ToFlightSeatResponse(s))
		}
	}

	return FlightResponse{
		ID:             f.ID,
		Airline:        ToAirlineResponse(f.Airline),
		Origin:         ToAirportResponse(f.Origin),
		Destination:    ToAirportResponse(f.Destination),
		DepartureTime:  f.DepartureTime,
		ArrivalTime:    f.ArrivalTime,
		FlightNumber:   f.FlightNumber,
		TotalSeats:     f.TotalSeats,
		AvailableSeats: available,
		TotalRows:      f.TotalRows,
		TotalColumns:   f.TotalColumns,
		FlightClasses:  classes,
		FlightSeats:    seats,
	}
}
