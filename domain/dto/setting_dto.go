package dto

type BaristaStatusRequest struct {
	IsAvailable *bool `json:"is_available" binding:"required"` // Pakai pointer agar nilai false bisa tervalidasi
}

type BaristaStatusResponse struct {
	IsAvailable bool   `json:"is_available"`
	StatusText  string `json:"status_text"`
}