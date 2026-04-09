package repository

import (
	"voca-store/internal/domain/models"
	"voca-store/internal/domain/repository"

	"gorm.io/gorm"
)

type dashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) repository.DashboardRepository {
	return &dashboardRepository{db: db}
}

func (r *dashboardRepository) CountUsers() (int64, error) {
	var count int64
	err := r.db.Model(&models.User{}).Count(&count).Error
	return count, err
}

func (r *dashboardRepository) CountProducts() (int64, error) {
	var count int64
	err := r.db.Model(&models.Product{}).Count(&count).Error
	return count, err
}

func (r *dashboardRepository) CountOrders() (int64, error) {
	var count int64
	err := r.db.Model(&models.Checkout{}).Count(&count).Error
	return count, err
}

func (r *dashboardRepository) SumRevenue() (float64, error) {
	var total float64
	err := r.db.
		Model(&models.Checkout{}).
		Select("COALESCE(SUM(total_price), 0)").
		Scan(&total).Error

	return total, err
}

func (r *dashboardRepository) CountPendingOrders() (int64, error) {
	var count int64
	err := r.db.Model(&models.Checkout{}).
		Where("status = ?", "pending").
		Count(&count).Error
	return count, err
}