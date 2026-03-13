package models

import (
	"time"

	"gorm.io/gorm"
)

type Airline struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"size:100;not null" json:"name"`
	Code         string         `gorm:"size:10;uniqueIndex;not null" json:"code"`
	LogoURL      string         `json:"logo_url"`
	LogoPublicID string         `json:"-"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type Airport struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Code      string         `gorm:"size:10;uniqueIndex;not null" json:"code"` // e.g., CGK
	Name      string         `gorm:"size:100;not null" json:"name"`
	City      string         `gorm:"size:100;not null" json:"city"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Flight struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	AirlineID      uint           `gorm:"not null;index" json:"airline_id"`
	Airline        Airline        `gorm:"foreignKey:AirlineID" json:"airline"`
	OriginID       uint           `gorm:"not null;index" json:"origin_id"`
	Origin         Airport        `gorm:"foreignKey:OriginID" json:"origin"`
	DestinationID  uint           `gorm:"not null;index" json:"destination_id"`
	Destination    Airport        `gorm:"foreignKey:DestinationID" json:"destination"`
	DepartureTime  time.Time      `gorm:"not null;index" json:"departure_time"`
	ArrivalTime    time.Time      `gorm:"not null" json:"arrival_time"`
	FlightNumber   string         `gorm:"size:20;not null" json:"flight_number"`
	TotalSeats     int            `gorm:"not null;default:0" json:"total_seats"`
	TotalRows      int            `gorm:"not null;default:0" json:"total_rows"`
	TotalColumns   int            `gorm:"not null;default:0" json:"total_columns"`
	AvailableSeats int            `json:"available_seats" gorm:"column:available_seats"`
	FlightClasses  []FlightClass  `gorm:"foreignKey:FlightID" json:"classes"`
	FlightSeats    []FlightSeat   `gorm:"foreignKey:FlightID" json:"flight_seats,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

type FlightClass struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	FlightID  uint           `gorm:"not null;index" json:"flight_id"`
	ClassType string         `gorm:"size:20;not null;index" json:"class_type"`
	Price     float64        `gorm:"not null" json:"price"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Seat is the master seat template table.
type Seat struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	SeatCode string `gorm:"size:10;uniqueIndex;not null" json:"seat_code"` // e.g. "1A"
}

// FlightSeat is the pivot table linking flights to seats.
type FlightSeat struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	FlightID      uint           `gorm:"not null;index;uniqueIndex:idx_flight_seat" json:"flight_id"`
	Flight        Flight         `gorm:"foreignKey:FlightID" json:"-"`
	SeatID        uint           `gorm:"not null;index;uniqueIndex:idx_flight_seat" json:"seat_id"`
	Seat          Seat           `gorm:"foreignKey:SeatID" json:"seat"`
	ClassType     string         `gorm:"size:20;not null;index" json:"class_type"`
	IsAvailable   bool           `gorm:"not null;default:true;index" json:"is_available"`
	LockedUntil   *time.Time     `json:"locked_until,omitempty"`
	TransactionID *uint          `gorm:"index" json:"transaction_id,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
