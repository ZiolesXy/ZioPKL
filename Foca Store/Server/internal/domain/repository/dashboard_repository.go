package repository

type DashboardRepository interface {
	CountUsers() (int64, error)
	CountProducts() (int64, error)
	CountOrders() (int64, error)
	SumRevenue() (float64, error)
	CountPendingOrders() (int64, error)
}