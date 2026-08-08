package dto

type OrderItemRequest struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
}

type CreateOrderRequest struct {
	Items       []OrderItemRequest `json:"items" binding:"required,min=1"`
	Discount    int                `json:"discount" binding:"min=0"`
	IsRedeem    bool               `json:"is_redeem"`
	RedeemToken string             `json:"redeem_token"` // Wajib diisi jika IsRedeem true
	SendEmailTo string             `json:"send_email_to"` // Opsional
	OrderVoucherCode string `json:"order_voucher_code"`
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