package repository

import (
	"context"
	"voca-plane/internal/domain/models"
	"voca-plane/pkg/helper"

	"gorm.io/gorm"
)

type adminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) AdminRepository {
	return &adminRepository{db: db}
}

func (r *adminRepository) GetDashboardStats(ctx context.Context) (*DashboardStats, error) {
	var stats DashboardStats

	r.db.WithContext(ctx).Model(&models.User{}).Count(&stats.TotalUsers)
	r.db.WithContext(ctx).Model(&models.Flight{}).Count(&stats.TotalFlights)
	r.db.WithContext(ctx).Model(&models.Transaction{}).Count(&stats.TotalTransactions)

	r.db.WithContext(ctx).Model(&models.Transaction{}).
		Where("payment_status = ?", "PAID").
		Select("COALESCE(SUM(total_price), 0)").Scan(&stats.TotalRevenue)

	r.db.WithContext(ctx).Model(&models.Transaction{}).
		Where("payment_status = ?", "PENDING").
		Count(&stats.PendingPayments)

	r.db.WithContext(ctx).Model(&models.Transaction{}).
		Where("payment_status = ?", "PAID").
		Count(&stats.CompletedBookings)

	return &stats, nil
}

func (r *adminRepository) GetAllUsers(ctx context.Context, page, limit int, sortBy, order string) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := r.db.WithContext(ctx).Unscoped().Model(&models.User{})
	query.Session(&gorm.Session{}).Count(&total)

	// Admin Users Whitelist
	allowedColumns := map[string]bool{
		"id":         true,
		"name":       true,
		"email":      true,
		"role":       true,
		"created_at": true,
		"updated_at": true,
	}

	query = helper.ApplySorting(query, sortBy, order, allowedColumns, "id ASC")

	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Find(&users).Error
	return users, total, err
}

func (r *adminRepository) UpdateUserRole(ctx context.Context, tx *gorm.DB, userID uint, role string) error {
	return tx.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Update("role", role).Error
}

func (r *adminRepository) GetAllTransactions(ctx context.Context, page, limit int, sortBy, order string) ([]models.Transaction, int64, error) {
	var transactions []models.Transaction
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Transaction{}).
		Preload("User").
		Preload("Flight").
		Preload("Flight.Airline").
		Preload("Flight.Origin").
		Preload("Flight.Destination").
		Preload("Items").
		Preload("Items.FlightClass")

	query.Session(&gorm.Session{}).Count(&total)

	// Transactions Whitelist
	allowedColumns := map[string]bool{
		"id":             true,
		"total_price":    true,
		"payment_status": true,
		"created_at":     true,
	}

	query = helper.ApplySorting(query, sortBy, order, allowedColumns, "id ASC")

	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Find(&transactions).Error
	return transactions, total, err
}