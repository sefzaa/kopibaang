package dto

type TopProductResponse struct {
	ProductName string `json:"product_name"`
	TotalSold   int    `json:"total_sold"`
}


type DashboardRequest struct {
	Filter    string `form:"filter"`     // today, yesterday, this_week, this_month, this_year, custom
	StartDate string `form:"start_date"` // Format: YYYY-MM-DD (Wajib jika filter=custom)
	EndDate   string `form:"end_date"`   // Format: YYYY-MM-DD (Wajib jika filter=custom)
}

type MenuStatResponse struct {
	ProductName string `json:"product_name"`
	TotalSold   int    `json:"total_sold"`
}

type DashboardResponse struct {
	TotalSales          int                `json:"total_sales"`
	TotalOrders         int                `json:"total_orders"`
	TotalPointsRedeemed int                `json:"total_points_redeemed"`
	MenuStats           []MenuStatResponse `json:"menu_stats"`
}