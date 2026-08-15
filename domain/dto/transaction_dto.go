package dto

import "time"

type OrderItemRequest struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
}

type CreateOrderRequest struct {
	Items            []OrderItemRequest `json:"items" binding:"required,min=1"`
	Discount         int                `json:"discount" binding:"min=0"`
	IsRedeem         bool               `json:"is_redeem"`
	RedeemToken      string             `json:"redeem_token"`       // Wajib diisi jika IsRedeem true
	SendEmailTo      string             `json:"send_email_to"`      // Opsional
	OrderVoucherCode string             `json:"order_voucher_code"` // Opsional
}

type CreateOrderResponse struct {
	OrderID        string `json:"order_id"`
	FinalAmount    int    `json:"final_amount"`
	EarnPointToken string `json:"earn_point_token,omitempty"` // QR Code untuk user scan
	ReceiptText    string `json:"receipt_text"`               // Teks untuk dilempar ke WA via Flutter
}

type ScanEarnPointRequest struct {
	EarnToken string `json:"earn_token" binding:"required"`
}

type RedeemQRResponse struct {
	RedeemToken string `json:"redeem_token"`
}

// === STRUKTUR BARU UNTUK FITUR HISTORY & FILTER ===

type OrderHistoryQueryRequest struct {
	Filter    string `form:"filter"`     // today, yesterday, this_week, this_month, custom
	StartDate string `form:"start_date"` // format: YYYY-MM-DD (Wajib jika filter='custom')
	EndDate   string `form:"end_date"`   // format: YYYY-MM-DD (Wajib jika filter='custom')
	Page      int    `form:"page"`       // Pagination
	Limit     int    `form:"limit"`      // Pagination
}

type PaginationRequest struct {
	Page  int `form:"page"`
	Limit int `form:"limit"`
}

type OrderItemDetailResponse struct {
	ProductID   string `json:"product_id"`
	Quantity    int    `json:"quantity"`
	PriceAtTime int    `json:"price_at_time"`
}

type OrderResponse struct {
	OrderID     string                    `json:"order_id"`
	UserID      string                    `json:"user_id,omitempty"`
	TotalAmount int                       `json:"total_amount"`
	Discount    int                       `json:"discount"`
	FinalAmount int                       `json:"final_amount"`
	IsRedeem    bool                      `json:"is_redeem"`
	Status      string                    `json:"status"`
	CreatedAt   time.Time                 `json:"created_at"`
	Items       []OrderItemDetailResponse `json:"items"`
}

type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalPages int   `json:"total_pages"`
	TotalItems int64 `json:"total_items"`
}

type OrderHistoryResponse struct {
	Orders []OrderResponse `json:"orders"`
	Meta   PaginationMeta  `json:"meta"`
}