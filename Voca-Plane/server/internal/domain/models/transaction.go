package models

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	ID            uint              `gorm:"primaryKey" json:"id"`
	Code          string            `gorm:"size:36;uniqueIndex;not null" json:"code"`
	UserID        uint              `gorm:"not null;index" json:"user_id"`
	User          User              `gorm:"foreignKey:UserID" json:"-"`
	FlightID      uint              `gorm:"not null" json:"flight_id"`
	Flight        Flight            `gorm:"foreignKey:FlightID" json:"flight"`
	TotalPrice    float64           `gorm:"not null" json:"total_price"`
	PaymentURL    string            `gorm:"size:255" json:"payment_url"`
	PaymentStatus string            `gorm:"size:20;default:'PENDING';index" json:"payment_status"`
	PromoCode     *string           `gorm:"size:50" json:"promo_code"`
	Discount      float64           `gorm:"default:0" json:"discount"`
	ExpiresAt     time.Time         `json:"expires_at"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
	DeletedAt     gorm.DeletedAt    `gorm:"index" json:"-"`
	Items         []TransactionItem `gorm:"foreignKey:TransactionID" json:"items"`
}

type TransactionItem struct {
	ID            uint        `gorm:"primaryKey" json:"id"`
	TransactionID uint        `gorm:"not null;index" json:"transaction_id"`
	PassengerName string      `gorm:"size:100;not null" json:"passenger_name"`
	Nationality   string      `gorm:"size:50" json:"nationality"`
	PassportNo    string      `gorm:"size:50" json:"passport_no"`
	SeatNumber    string      `gorm:"size:10;not null" json:"seat_number"`
	FlightSeatID  uint        `gorm:"not null;index" json:"flight_seat_id"`
	FlightClassID uint        `gorm:"not null;index" json:"flight_class_id"`
	FlightClass   FlightClass `gorm:"foreignKey:FlightClassID" json:"flight_class"`
	Price         float64     `gorm:"not null" json:"price"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

type PromoCode struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Code      string         `gorm:"size:50;uniqueIndex;not null" json:"code"`
	Discount  float64        `gorm:"not null" json:"discount"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
