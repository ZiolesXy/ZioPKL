package response

import (
	"time"
	"voca-plane/internal/domain/models"
)

type TransactionItemResponse struct {
	PassengerName string `json:"passenger_name"`
	Nationality   string `json:"nationality"`
	PassportNo    string `json:"passport_no"`
	SeatNumber    string `json:"seat_number"`
	ClassName     string `json:"class_name"`
	Price         float64 `json:"price"`
}

type SimpleFlightResponse struct {
	FlightNumber  string          `json:"flight_number"`
	DepartureTime time.Time       `json:"departure_time"`
	ArrivalTime   time.Time       `json:"arrival_time"`
	Airline       AirlineResponse `json:"airline"`
	Origin        AirportResponse `json:"origin"`
	Destination   AirportResponse `json:"destination"`
}

type TransactionResponse struct {
	ID            uint                      `json:"id"`
	Code          string                    `json:"code"`
	Flight        SimpleFlightResponse      `json:"flight"`
	TotalPrice    float64                   `json:"total_price"`
	PaymentStatus string                    `json:"payment_status"`
	PaymentURL    string                    `json:"payment_url"`
	PromoCode     *string                   `json:"promo_code,omitempty"`
	Discount      float64                   `json:"discount"`
	ExpiresAt     time.Time                 `json:"expires_at"`
	CreatedAt     time.Time                 `json:"created_at"`
	Items         []TransactionItemResponse `json:"transactions_passangers,omitempty"`
}

func ToTransactionItemResponse(p models.TransactionItem) TransactionItemResponse {
	return TransactionItemResponse{
		PassengerName: p.PassengerName,
		Nationality:   p.Nationality,
		PassportNo:    p.PassportNo,
		SeatNumber:    p.SeatNumber,
		ClassName:     p.FlightClass.ClassType,
		Price:         p.Price,
	}
}

func ToSimpleFlightResponse(f models.Flight) SimpleFlightResponse {
	return SimpleFlightResponse{
		FlightNumber:  f.FlightNumber,
		DepartureTime: f.DepartureTime,
		ArrivalTime:   f.ArrivalTime,
		Airline:       ToAirlineResponse(f.Airline),
		Origin:        ToAirportResponse(f.Origin),
		Destination:   ToAirportResponse(f.Destination),
	}
}

func ToTransactionResponse(t models.Transaction) TransactionResponse {
	var items []TransactionItemResponse
	for _, p := range t.Items {
		items = append(items, ToTransactionItemResponse(p))
	}

	return TransactionResponse{
		ID:            t.ID,
		Code:          t.Code,
		Flight:        ToSimpleFlightResponse(t.Flight),
		TotalPrice:    t.TotalPrice,
		PaymentStatus: t.PaymentStatus,
		PaymentURL:    t.PaymentURL,
		PromoCode:     t.PromoCode,
		Discount:      t.Discount,
		ExpiresAt:     t.ExpiresAt,
		CreatedAt:     t.CreatedAt,
		Items:         items,
	}
}
