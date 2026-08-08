package dto

type RawMaterialRequest struct {
	Name     string  `json:"name" binding:"required"`
	Quantity float64 `json:"quantity" binding:"required,min=0"`
	Unit     string  `json:"unit" binding:"required"`
	Price    int     `json:"price" binding:"required,min=0"`
	Source   string  `json:"source"`
}

type RawMaterialResponse struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Quantity  float64 `json:"quantity"`
	Unit      string  `json:"unit"`
	Price     int     `json:"price"`
	Source    string  `json:"source"`
	CreatedAt string  `json:"created_at"`
}