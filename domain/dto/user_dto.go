package dto

type UpdateProfileRequest struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
}

type UserProfileResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"` // SUDAH DITAMBAHKAN
	Points   int    `json:"points"`
	Role     string `json:"role"`
}