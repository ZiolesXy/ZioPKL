package service

import (
	"voca-store/internal/domain/dto/response"
	"voca-store/internal/domain/repository"
)

type AdminService interface {
	GetDashboardData() (*response.DashboardData, error)
}

type adminService struct {
	dashboardRepo repository.DashboardRepository
}

func NewAdminService(dashboardRepo repository.DashboardRepository) AdminService {
	return &adminService{dashboardRepo}
}

func (s *adminService) GetDashboardData() (*response.DashboardData, error) {
	users, err := s.dashboardRepo.CountUsers()
	if err != nil {
		return nil, err
	}

	category, err := s.dashboardRepo.CountCategory()
	if err != nil {
		return nil, err
	}

	products, err := s.dashboardRepo.CountProducts()
	if err != nil {
		return nil, err
	}

	coupons, err := s.dashboardRepo.CountCoupons()
	if err != nil {
		return nil, err
	}

	orders, err := s.dashboardRepo.CountOrders()
	if err != nil {
		return nil, err
	}

	revenue, err := s.dashboardRepo.SumRevenue()
	if err != nil {
		return nil, err
	}

	pending, err := s.dashboardRepo.CountPendingOrders()
	if err != nil {
		return nil, err
	}

	accepted, err := s.dashboardRepo.CountAcceptedOrders()
	if err != nil {
		return nil, err
	}

	decline, err := s.dashboardRepo.CountDeclineOrders()
	if err != nil {
		return nil, err
	}

	return &response.DashboardData{
		TotalUsers:    users,
		TotalCategory: category,
		TotalProducts: products,
		TotalCoupons: coupons,
		TotalOrders:   orders,
		TotalRevenue:  revenue,
		PendingOrders: pending,
		AcceptedOrders: accepted,
		DeclineOrders: decline,
	}, nil
}