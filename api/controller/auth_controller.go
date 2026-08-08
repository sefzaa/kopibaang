package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"kopibang/domain/dto"
	"kopibang/usecase"
)

type AuthController struct {
	authUsecase *usecase.AuthUsecase
}

func NewAuthController(authUsecase *usecase.AuthUsecase) *AuthController {
	return &AuthController{authUsecase}
}

// LoginCustomer godoc
// @Summary Customer Login
// @Description Authenticate a customer and return tokens
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login Credentials"
// @Success 200 {object} dto.BaseResponse{data=dto.TokenResponse}
// @Failure 401 {object} dto.BaseResponse
// @Router /api/v1/auth/login [post]
func (c *AuthController) LoginCustomer(ctx *gin.Context) {
	var req dto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(http.StatusBadRequest, "Invalid request body", err.Error()))
		return
	}

	res, err := c.authUsecase.Login(ctx.Request.Context(), req, "customer")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse(http.StatusUnauthorized, err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse(http.StatusOK, "Login success", res))
}

// LoginAdmin godoc
// @Summary Admin (Barista) Login
// @Description Authenticate a barista and return tokens
// @Tags Admin Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login Credentials"
// @Success 200 {object} dto.BaseResponse{data=dto.TokenResponse}
// @Failure 401 {object} dto.BaseResponse
// @Router /api/v1/admin/auth/login [post]
func (c *AuthController) LoginAdmin(ctx *gin.Context) {
	var req dto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(http.StatusBadRequest, "Invalid request body", err.Error()))
		return
	}

	res, err := c.authUsecase.Login(ctx.Request.Context(), req, "barista")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse(http.StatusUnauthorized, err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse(http.StatusOK, "Login success", res))
}

// RefreshToken godoc
// @Summary Refresh Access Token
// @Description Get a new access token and refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh Token"
// @Success 200 {object} dto.BaseResponse{data=dto.TokenResponse}
// @Failure 401 {object} dto.BaseResponse
// @Router /api/v1/auth/refresh [post]
func (c *AuthController) RefreshToken(ctx *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(http.StatusBadRequest, "Invalid request body", err.Error()))
		return
	}

	res, err := c.authUsecase.RefreshToken(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse(http.StatusUnauthorized, err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse(http.StatusOK, "Token refreshed successfully", res))
}

// RegisterCustomer godoc
// @Summary Register a new customer
// @Description Register and auto-login customer
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Registration Info"
// @Success 201 {object} dto.BaseResponse{data=dto.TokenResponse}
// @Failure 400 {object} dto.BaseResponse
// @Router /api/v1/auth/register [post]
func (c *AuthController) RegisterCustomer(ctx *gin.Context) {
	var req dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(http.StatusBadRequest, "Invalid request body or passwords do not match", err.Error()))
		return
	}

	res, err := c.authUsecase.Register(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusCreated, dto.SuccessResponse(http.StatusCreated, "Registration and auto-login success", res))
}

// ForgotPassword godoc
// @Summary Request Password Reset
// @Description Sends OTP to registered email
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.ForgotPasswordRequest true "Email address"
// @Success 200 {object} dto.BaseResponse
// @Failure 404 {object} dto.BaseResponse
// @Router /api/v1/auth/forgot-password [post]
func (c *AuthController) ForgotPassword(ctx *gin.Context) {
	var req dto.ForgotPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(http.StatusBadRequest, "Invalid request body", err.Error()))
		return
	}

	err := c.authUsecase.ForgotPassword(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.ErrorResponse(http.StatusNotFound, err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse(http.StatusOK, "OTP has been sent to your email", nil))
}



// Logout godoc
// @Summary Logout User/Admin
// @Description Invalidate the refresh token in Redis
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.LogoutRequest true "Refresh Token"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseResponse
// @Router /api/v1/auth/logout [post]
func (c *AuthController) Logout(ctx *gin.Context) {
	var req dto.LogoutRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(http.StatusBadRequest, "Invalid request body", err.Error()))
		return
	}

	err := c.authUsecase.Logout(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse(http.StatusOK, "Logout successful", nil))
}


// VerifyOTP godoc
// @Summary Verify OTP
// @Description Verify OTP sent to email and get reset token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.VerifyOTPRequest true "Email and OTP"
// @Success 200 {object} dto.BaseResponse{data=dto.VerifyOTPResponse}
// @Failure 400 {object} dto.BaseResponse
// @Router /api/v1/auth/verify-otp [post]
func (c *AuthController) VerifyOTP(ctx *gin.Context) {
	var req dto.VerifyOTPRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(http.StatusBadRequest, "Invalid request body", err.Error()))
		return
	}

	res, err := c.authUsecase.VerifyOTP(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse(http.StatusOK, "OTP verified successfully", res))
}

// ResetPassword godoc
// @Summary Reset Password
// @Description Submit reset token and new password to reset
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.ResetPasswordRequest true "Email, Reset Token, and New Password"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseResponse
// @Router /api/v1/auth/reset-password [post]
func (c *AuthController) ResetPassword(ctx *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(http.StatusBadRequest, "Invalid request body or passwords do not match", err.Error()))
		return
	}

	err := c.authUsecase.ResetPassword(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse(http.StatusOK, "Password reset successful", nil))
}