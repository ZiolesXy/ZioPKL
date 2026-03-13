package repository

import (
	"context"
	"time"
	"voca-plane/internal/domain/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, id uint) (*models.User, error)

	Delete(ctx context.Context, id uint) error
	Restore(ctx context.Context, id uint) error
	Ban(ctx context.Context, id uint, reason string) error
	Unban(ctx context.Context, id uint) error
}

type FlightRepository interface {
	Search(ctx context.Context, origin, destination, date, classType string, page, limit int) ([]models.Flight, int64, error)
	GetByID(ctx context.Context, id uint) (*models.Flight, error)
	GetClassByID(ctx context.Context, id uint) (*models.FlightClass, error)
	GetAll(ctx context.Context, page, limit int, sortBy, order string) ([]models.Flight, int64, error)
	GetAllFull(ctx context.Context) ([]models.Flight, error)
	Create(ctx context.Context, tx *gorm.DB, flight *models.Flight) error
	Update(ctx context.Context, tx *gorm.DB, flight *models.Flight) error
	Delete(ctx context.Context, tx *gorm.DB, id uint) error
	GetFlightWithClasses(ctx context.Context, tx *gorm.DB, id uint) (*models.Flight, error)
	GetFlightWithRelations(ctx context.Context, tx *gorm.DB, id uint) (*models.Flight, error)
	BulkCreateFlightSeats(ctx context.Context, tx *gorm.DB, seats []models.FlightSeat) error
	DeleteFlightSeatsByFlightID(ctx context.Context, tx *gorm.DB, flightID uint) error
	CreateClass(ctx context.Context, tx *gorm.DB, class *models.FlightClass) error
	GetOrCreateSeats(ctx context.Context, tx *gorm.DB, codes []string) ([]models.Seat, error)
	GetAvailableSeats(ctx context.Context, flightID uint, classType string) ([]models.FlightSeat, error)
	GetFlightSeatsByIDs(ctx context.Context, tx *gorm.DB, ids []uint) ([]models.FlightSeat, error)
	GetFlightSeatsByCodes(ctx context.Context, tx *gorm.DB, flightID uint, codes []string) ([]models.FlightSeat, error)
	LockSeats(ctx context.Context, tx *gorm.DB, seatIDs []uint, transactionID uint, until time.Time) error
	UnlockExpiredSeats(ctx context.Context) error
	ReleaseSeats(ctx context.Context, tx *gorm.DB, transactionID uint) error
	FinalizeSeats(ctx context.Context, tx *gorm.DB, transactionID uint) error
}

type AirlineRepository interface {
	GetAll(ctx context.Context, page, limit int, sortBy, order string) ([]models.Airline, int64, error)
	GetByID(ctx context.Context, id uint) (*models.Airline, error)
	Create(ctx context.Context, tx *gorm.DB, airline *models.Airline) error
	Update(ctx context.Context, tx *gorm.DB, airline *models.Airline) error
	Delete(ctx context.Context, tx *gorm.DB, id uint) error
}

type AirportRepository interface {
	GetAll(ctx context.Context, page, limit int, sortBy, order string) ([]models.Airport, int64, error)
	GetByID(ctx context.Context, id uint) (*models.Airport, error)
	Create(ctx context.Context, tx *gorm.DB, airport *models.Airport) error
	Update(ctx context.Context, tx *gorm.DB, airport *models.Airport) error
	Delete(ctx context.Context, tx *gorm.DB, id uint) error
}

type TransactionRepository interface {
	Create(ctx context.Context, tx *gorm.DB, t *models.Transaction) error
	GetByCode(ctx context.Context, code string) (*models.Transaction, error)
	GetByUserID(ctx context.Context, userID uint, page, limit int) ([]models.Transaction, int64, error)
	GetByUserIDAll(ctx context.Context, userID uint) ([]models.Transaction, error)
	UpdatePaymentStatus(ctx context.Context, tx *gorm.DB, id uint, status string) error
	Delete(ctx context.Context, tx *gorm.DB, code string) error
	CreateTransactionItems(ctx context.Context, tx *gorm.DB, items []models.TransactionItem) error
	UpdatePaymentURL(ctx context.Context, code string, url string) error
}

type PromoRepository interface {
	GetByID(ctx context.Context, id uint) (*models.PromoCode, error)
	GetByCode(ctx context.Context, code string) (*models.PromoCode, error)
	GetAll(ctx context.Context, page, limit int, sortBy, order string) ([]models.PromoCode, int64, error)
	Create(ctx context.Context, tx *gorm.DB, promo *models.PromoCode) error
	Update(ctx context.Context, tx *gorm.DB, promo *models.PromoCode) error
	Delete(ctx context.Context, tx *gorm.DB, id uint) error
}

type DashboardStats struct {
	TotalUsers        int64   `json:"total_users"`
	TotalFlights      int64   `json:"total_flights"`
	TotalTransactions int64   `json:"total_transactions"`
	TotalRevenue      float64 `json:"total_revenue"`
	PendingPayments   int64   `json:"pending_payments"`
	CompletedBookings int64   `json:"completed_bookings"`
}

type AdminRepository interface {
	GetDashboardStats(ctx context.Context) (*DashboardStats, error)
	GetAllUsers(ctx context.Context, page, limit int, sortBy, order string) ([]models.User, int64, error)
	UpdateUserRole(ctx context.Context, tx *gorm.DB, userID uint, role string) error
	GetAllTransactions(ctx context.Context, page, limit int, sortBy, order string) ([]models.Transaction, int64, error)
}

type SystemRepository interface {
	ResetDatabase(ctx context.Context) error
}