package dto

type DashboardStatsResponse struct {
	TotalUsers       int64 `json:"total_users"`
	TotalGames       int64 `json:"total_games"`
	TotalActiveChats int64 `json:"total_active_chats"`
}
