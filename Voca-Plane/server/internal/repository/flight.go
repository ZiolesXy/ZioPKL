package repository

import (
	"context"
	"time"
	"voca-plane/internal/domain/models"
	"voca-plane/pkg/helper"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type flightRepository struct {
	db *gorm.DB
}

func NewFlightRepository(db *gorm.DB) FlightRepository {
	return &flightRepository{db: db}
}

func (r *flightRepository) Search(ctx context.Context, origin, destination, date, classType string, page, limit int) ([]models.Flight, int64, error) {
	var flights []models.Flight
	var total int64

	layout := "2006-01-02"
	parseDate, err := time.Parse(layout, date)
	if err != nil {
		return nil, 0, err
	}
	startOfDay := time.Date(parseDate.Year(), parseDate.Month(), parseDate.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24*time.Hour - time.Second)

	query := r.db.WithContext(ctx).Model(&models.Flight{}).
		Joins("JOIN airports AS origin ON flights.origin_id = origin.id").
		Joins("JOIN airports AS dest ON flights.destination_id = dest.id").
		Preload("Airline").
		Preload("Origin").
		Preload("Destination").
		Preload("FlightClasses").
		Preload("FlightSeats.Seat").
		Where("origin.name LIKE %?% AND dest.name LIKE ?", "%"+origin+"%", "%"+destination+"%").
		Where("flights.departure_time BETWEEN ? AND ?", startOfDay, endOfDay)
	
	if classType != "" {
		query = query.Joins("JOIN flight_classes ON flight_classes.flight_id = flights.id").
			Where("flight_classes.class_type = ?", classType)
	}

	query.Session(&gorm.Session{}).Count(&total)
	offset := (page - 1) * limit
	err = query.Offset(offset).Limit(limit).Find(&flights).Error
	return flights, total, err
}

func (r *flightRepository) GetByID(ctx context.Context, id uint) (*models.Flight, error) {
	var flight models.Flight
	err := r.db.WithContext(ctx).
		Preload("Airline").
		Preload("Origin").
		Preload("Destination").
		Preload("FlightClasses").
		Preload("FlightSeats.Seat").
		First(&flight, id).Error
	return &flight, err
}

func (r *flightRepository) GetClassByID(ctx context.Context, id uint) (*models.FlightClass, error) {
	var class models.FlightClass
	err := r.db.WithContext(ctx).First(&class, id).Error
	return &class, err
}

func (r *flightRepository) GetAll(ctx context.Context, page, limit int, sortBy, order string) ([]models.Flight, int64, error) {
	type flightScan struct {
		models.Flight
		AvailableSeats int `gorm:"column:available_seats"`
	}

	var scannedFlights []flightScan
	var total int64

	available := true
	now := time.Now()
	query := r.db.WithContext(ctx).Model(&models.Flight{}).
		Select(`
			flights.*,
			(
				SELECT COUNT(*)
				FROM flight_seats
				WHERE flight_seats.flight_id = flights.id
				AND flight_seats.is_available = ?
				AND (flight_seats.locked_until IS NULL OR flight_seats.locked_until < ?)
			) AS available_seats
		`, available, now).
		Preload("Airline").
		Preload("Origin").
		Preload("Destination").
		Preload("FlightClasses")
	
	query.Session(&gorm.Session{}).Count(&total)

	// Flights Whitelist
	allowedColumns := map[string]bool{
		"id":             true,
		"departure_time": true,
		"arrival_time":   true,
		"total_seats":    true,
	}

	query = helper.ApplySorting(query, sortBy, order, allowedColumns, "id ASC")

	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Find(&scannedFlights).Error
	if err != nil {
		return nil, 0, err
	}

	flights := make([]models.Flight, len(scannedFlights))
	for i, s := range scannedFlights {
		flights[i] = s.Flight
		flights[i].AvailableSeats = s.AvailableSeats
	}

	return flights, total, nil
}

func(r *flightRepository) GetAllFull(ctx context.Context) ([]models.Flight, error) {
	type flightScan struct {
		models.Flight
		AvailableSeats int `gorm:"column:available_seats"`
	}

	var scannedFlights []flightScan

	available := true
	now := time.Now()
	err := r.db.WithContext(ctx).
		Model(&models.Flight{}).
		Select(`
			flights.*,
			(
				SELECT COUNT(*)
				FROM flight_seats
				WHERE flight_seats.flight_id = flights.id
				AND flight_seats.is_available = ?
				AND (flight_seats.locked_until IS NULL OR flight_seats.locked_until < ?)
			) AS available_seats
		`, available, now).
		Preload("Airline").
		Preload("Origin").
		Preload("Destination").
		Preload("FlightClasses").
		Order("departure_time ASC").
		Find(&scannedFlights).Error

	if err != nil {
		return nil, err
	}

	flights := make([]models.Flight, len(scannedFlights))
	for i, s := range scannedFlights {
		flights[i] = s.Flight
		flights[i].AvailableSeats = s.AvailableSeats
	}

	return flights, nil
}

func (r *flightRepository) Create(ctx context.Context, tx *gorm.DB, flight *models.Flight) error {
	if err := tx.WithContext(ctx).Create(flight).Error; err != nil {
		return err
	}

	return tx.WithContext(ctx).
		Preload("Airline").
		Preload("Origin").
		Preload("Destination").
		Preload("FlightClasses").
		Preload("FlightSeats.Seat").
		First(flight, flight.ID).Error
}

func (r *flightRepository) Update(ctx context.Context, tx *gorm.DB, flight *models.Flight) error {
	if err := tx.WithContext(ctx).Save(flight).Error; err != nil {
		return err
	}

	return tx.WithContext(ctx).
		Preload("Airline").
		Preload("Origin").
		Preload("Destination").
		Preload("FlightClasses").
		Preload("FlightSeats.Seat").
		First(flight, flight.ID).Error
}

func (r *flightRepository) Delete(ctx context.Context, tx *gorm.DB, id uint) error {
	return tx.WithContext(ctx).Delete(&models.Flight{}, id).Error
}

func (r *flightRepository) GetFlightWithClasses(ctx context.Context, tx *gorm.DB, id uint) (*models.Flight, error) {
	var flight models.Flight

	err := tx.WithContext(ctx).
		Preload("FlightClasses").
		First(&flight, id).Error

	if err != nil {
		return nil, err
	}

	return &flight, nil
}

func (r *flightRepository) GetFlightWithRelations(ctx context.Context, tx *gorm.DB, id uint) (*models.Flight, error) {
	var flight models.Flight

	err := tx.WithContext(ctx).
		Preload("Airline").
		Preload("Origin").
		Preload("Destination").
		Preload("FlightClasses").
		Preload("FlightSeats.Seat").
		First(&flight, id).Error

	if err != nil {
		return nil, err
	}

	return &flight, nil
}

func (r *flightRepository) BulkCreateFlightSeats(ctx context.Context, tx *gorm.DB, seats []models.FlightSeat) error {
	return tx.WithContext(ctx).Create(&seats).Error
}

func (r *flightRepository) DeleteFlightSeatsByFlightID(ctx context.Context, tx *gorm.DB, flightID uint) error {
	return tx.WithContext(ctx).
		Unscoped().
		Where("flight_id = ?", flightID).
		Delete(&models.FlightSeat{}).Error
}

func (r *flightRepository) CreateClass(ctx context.Context, tx *gorm.DB, class *models.FlightClass) error {
	return tx.WithContext(ctx).Create(class).Error
}

func (r *flightRepository) GetOrCreateSeats(ctx context.Context, tx *gorm.DB, codes []string) ([]models.Seat, error) {
	// Build seat models
	seatModels := make([]models.Seat, len(codes))
	for i, code := range codes {
		seatModels[i] = models.Seat{SeatCode: code}
	}

	// Upsert: insert if not exists, do nothing on conflict
	if err := tx.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "seat_code"}},
			DoNothing: true,
		}).
		Create(&seatModels).Error; err != nil {
		return nil, err
	}

	// Fetch all seats by codes to get their IDs (including pre-existing ones)
	var seats []models.Seat
	if err := tx.WithContext(ctx).Where("seat_code IN ?", codes).Find(&seats).Error; err != nil {
		return nil, err
	}

	return seats, nil
}

func (r *flightRepository) GetAvailableSeats(ctx context.Context, flightID uint, classType string) ([]models.FlightSeat, error) {
	var seats []models.FlightSeat

	query := r.db.WithContext(ctx).
		Preload("Seat").
		Where("flight_id = ? AND is_available = ?", flightID, true).
		Where("locked_until IS NULL OR locked_until < ?", time.Now())

	if classType != "" {
		query = query.Where("class_type = ?", classType)
	}

	err := query.Find(&seats).Error
	return seats, err
}

func (r *flightRepository) GetFlightSeatsByIDs(ctx context.Context, tx *gorm.DB, ids []uint) ([]models.FlightSeat, error) {
	var seats []models.FlightSeat
	err := tx.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Preload("Seat").
		Where("id IN ?", ids).
		Find(&seats).Error
	return seats, err
}

// GetFlightSeatsByCodes fetches specific pivot rows by seat codes with pessimistic locking (FOR UPDATE).
func (r *flightRepository) GetFlightSeatsByCodes(ctx context.Context, tx *gorm.DB, flightID uint, codes []string) ([]models.FlightSeat, error) {
	var seats []models.FlightSeat
	err := tx.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Preload("Seat").
		Joins("JOIN seats ON seats.id = flight_seats.seat_id").
		Where("flight_seats.flight_id = ? AND seats.seat_code IN ?", flightID, codes).
		Find(&seats).Error
	return seats, err
}

func (r *flightRepository) LockSeats(ctx context.Context, tx *gorm.DB, seatIDs []uint, transactionID uint, until time.Time) error {
	return tx.WithContext(ctx).
		Model(&models.FlightSeat{}).
		Where("id IN ?", seatIDs).
		Updates(map[string]interface{}{
			"is_available":   false,
			"locked_until":   until,
			"transaction_id": transactionID,
		}).Error
}

func (r *flightRepository) UnlockExpiredSeats(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Model(&models.FlightSeat{}).
		Where("locked_until IS NOT NULL AND locked_until < ? AND is_available = ?", time.Now(), false).
		Updates(map[string]interface{}{
			"is_available":   true,
			"locked_until":   nil,
			"transaction_id": nil,
		}).Error
}

func (r *flightRepository) ReleaseSeats(ctx context.Context, tx *gorm.DB, transactionID uint) error {
	return tx.WithContext(ctx).
		Model(&models.FlightSeat{}).
		Where("transaction_id = ?", transactionID).
		Updates(map[string]interface{}{
			"is_available":   true,
			"locked_until":   nil,
			"transaction_id": nil,
		}).Error
}

func (r *flightRepository) FinalizeSeats(ctx context.Context, tx *gorm.DB, transactionID uint) error {
	return tx.WithContext(ctx).
		Model(&models.FlightSeat{}).
		Where("transaction_id = ?", transactionID).
		Updates(map[string]interface{}{
			"locked_until": nil,
		}).Error
}