package request

type CreateTransactionRequest struct {
	FlightID   uint               `json:"flight_id" binding:"required"`
	Passengers []PassengerRequest `json:"passengers" binding:"required,dive"`
	PromoCode  *string            `json:"promo_code"`
}

type PassengerRequest struct {
	FullName    string `json:"full_name" binding:"required"`
	Nationality string `json:"nationality"`
	PassportNo  string `json:"passport_no" binding:"required"`
	SeatNumber  string `json:"seat_number" binding:"required"`
}