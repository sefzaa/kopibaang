package dto

type IngredientRequest struct {
	Name     string `json:"name" binding:"required"`
	Grammage string `json:"grammage" binding:"required"`
}

type IngredientResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Grammage string `json:"grammage"`
}

type ProductRequest struct {
	Name        string              `json:"name" binding:"required"`
	Description string              `json:"description"`
	Price       int                 `json:"price" binding:"required,min=0"`
	Discount    int                 `json:"discount" binding:"min=0"`
	VoucherID   *string             `json:"voucher_id"` 
	Volume      string              `json:"volume" binding:"required"`
	IsActive    *bool               `json:"is_active" binding:"required"`
	Ingredients []IngredientRequest `json:"ingredients" binding:"required,dive"`
	ImageURLs   []string            `json:"image_urls"` // Hanya pakai versi Array
}

type ProductResponse struct {
	ID                 string               `json:"id"`
	Name               string               `json:"name"`
	Description        string               `json:"description"`
	Price              int                  `json:"price"`
	Discount           int                  `json:"discount"`
	VoucherID          *string              `json:"voucher_id,omitempty"`
	VoucherCode        string               `json:"voucher_code,omitempty"`
	VoucherDiscount    int                  `json:"voucher_discount"` 
	FinalPrice         int                  `json:"final_price"`      
	Volume             string               `json:"volume"`
	IsActive           bool                 `json:"is_active"`
	Ingredients        []IngredientResponse `json:"ingredients"`
	ImageURLs          []string             `json:"image_urls"` // Hanya pakai versi Array
}