package dto

type DashboardStatsResponse struct {
	TotalUsers       int64 `json:"total_users"`
	TotalDevelopers  int64 `json:"total_developers"`
	TotalGames       int64 `json:"total_games"`
	TotalActiveChats int64 `json:"total_active_chats"`
}
