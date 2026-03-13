package request

type SearchFlightRequest struct {
	Origin      string `form:"origin" binding:"required"`
	Destination string `form:"destination" binding:"required"`
	Date        string `form:"date" binding:"required"`
	ClassType   string `form:"class_type"`
	Page        int    `form:"page"`
	Limit       int    `form:"limit"`
}

type CreateFlightRequest struct {
	AirlineID     uint                `json:"airline_id" binding:"required"`
	OriginID      uint                `json:"origin_id" binding:"required"`
	DestinationID uint                `json:"destination_id" binding:"required"`
	DepartureTime string              `json:"departure_time" binding:"required"`
	ArrivalTime   string              `json:"arrival_time" binding:"required"`
	FlightNumber  string              `json:"flight_number" binding:"required"`
	TotalSeats    int                 `json:"total_seats" binding:"required,min=1"`
	TotalRows     int                 `json:"total_rows" binding:"required,min=1,max=26"`
	TotalColumns  int                 `json:"total_columns" binding:"required,min=1,max=10"`
	ClassCount    int                 `json:"class_count" binding:"required,min=1,max=3"`
	ClassPrices   []ClassPriceRequest `json:"class_prices" binding:"required,dive"`
}

type ClassPriceRequest struct {
	ClassType string  `json:"class_type" binding:"required"`
	Price     float64 `json:"price" binding:"required,min=0"`
}

type UpdateFlightRequest struct {
	AirlineID     *uint   `json:"airline_id,omitempty"`
	OriginID      *uint   `json:"origin_id,omitempty"`
	DestinationID *uint   `json:"destination_id,omitempty"`
	DepartureTime *string `json:"departure_time,omitempty"`
	ArrivalTime   *string `json:"arrival_time,omitempty"`
	FlightNumber  *string `json:"flight_number,omitempty"`
	TotalSeats    *int    `json:"total_seats,omitempty"`
	TotalRows     *int    `json:"total_rows,omitempty"`
	TotalColumns  *int    `json:"total_columns,omitempty"`
}