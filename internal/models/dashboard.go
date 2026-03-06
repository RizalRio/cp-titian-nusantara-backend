package models

// DashboardStatsResponse adalah DTO untuk statistik dashboard admin
type DashboardStatsResponse struct {
	TotalPosts     int64 `json:"total_posts"`
	TotalServices  int64 `json:"total_services"`
	TotalPortfolios int64 `json:"total_portfolios"`
	UnreadMessages int64 `json:"unread_messages"`
}