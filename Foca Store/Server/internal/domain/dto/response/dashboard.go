package response

type DashboardData struct {
	TotalUsers     int64
	TotalCategory  int64
	TotalProducts  int64
	TotalCoupons   int64
	TotalOrders    int64
	TotalRevenue   float64
	PendingOrders  int64
	AcceptedOrders int64
	DeclineOrders  int64
}

type DashboardResponse struct {
	TotalUsers     int64   `json:"total_users"`
	TotalCategory  int64   `json:"total_category"`
	TotalProducts  int64   `json:"total_products"`
	TotalCoupons   int64   `json:"total_coupons"`
	TotalOrders    int64   `json:"total_orders"`
	TotalRevenue   float64 `json:"total_revenue"`
	PendingOrders  int64   `json:"pending_orders"`
	AcceptedOrders int64   `json:"accepted_orders"`
	DeclineOrders  int64   `json:"decline_orders"`
}

func BuildDashboardResponse(data DashboardData) DashboardResponse {
	return DashboardResponse{
		TotalUsers:     data.TotalUsers,
		TotalCategory:  data.TotalCategory,
		TotalProducts:  data.TotalProducts,
		TotalCoupons:   data.TotalCoupons,
		TotalOrders:    data.TotalOrders,
		TotalRevenue:   data.TotalRevenue,
		PendingOrders:  data.PendingOrders,
		AcceptedOrders: data.AcceptedOrders,
		DeclineOrders:  data.DeclineOrders,
	}
}
