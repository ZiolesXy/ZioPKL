package response

type DashboardData struct {
	TotalUsers    int64
	TotalProducts int64
	TotalOrders   int64
	TotalRevenue  float64
	PendingOrders int64
}

type DashboardResponse struct {
	TotalUsers    int64   `json:"total_users"`
	TotalProducts int64   `json:"total_products"`
	TotalOrders   int64   `json:"total_orders"`
	TotalRevenue  float64 `json:"total_revenue"`
	PendingOrders int64   `json:"pending_orders"`
}

func BuildDashboardResponse(data DashboardData) DashboardResponse {
	return DashboardResponse{
		TotalUsers: data.TotalUsers,
		TotalProducts: data.TotalProducts,
		TotalOrders: data.TotalOrders,
		TotalRevenue: data.TotalRevenue,
		PendingOrders: data.PendingOrders,
	}
}