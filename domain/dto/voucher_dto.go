package dto

type VoucherRequest struct {
	Code           string `json:"code" binding:"required"`
	Type           string `json:"type" binding:"required,oneof=menu_promo cart_discount"` // TAMBAHAN: Validasi ketat enum
	DiscountAmount int    `json:"discount_amount" binding:"required,min=1"`
	MinPurchase    int    `json:"min_purchase" binding:"min=0"`
	IsActive       *bool  `json:"is_active" binding:"required"`
	StartDate      string `json:"start_date" binding:"required"` // Format: YYYY-MM-DD
	EndDate        string `json:"end_date" binding:"required"`   // Format: YYYY-MM-DD
}

type VoucherResponse struct {
	ID             string `json:"id"`
	Code           string `json:"code"`
	Type           string `json:"type"` // TAMBAHAN
	DiscountAmount int    `json:"discount_amount"`
	MinPurchase    int    `json:"min_purchase"`
	IsActive       bool   `json:"is_active"`
	StartDate      string `json:"start_date"`
	EndDate        string `json:"end_date"`
}