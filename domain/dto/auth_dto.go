package dto

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Name            string `json:"name" binding:"required"`
	Username        string `json:"username" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=6"`
	RewritePassword string `json:"rewrite_password" binding:"required,eqfield=Password"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Role         string `json:"role"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type VerifyOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required,len=6"`
}

type VerifyOTPResponse struct {
	ResetToken string `json:"reset_token"`
}

type ResetPasswordRequest struct {
	Email           string `json:"email" binding:"required,email"`
	ResetToken      string `json:"reset_token" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
	RewritePassword string `json:"rewrite_password" binding:"required,eqfield=NewPassword"`
}