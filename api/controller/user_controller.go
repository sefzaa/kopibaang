package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"kopibang/domain/dto"
	"kopibang/usecase"
)

type UserController struct {
	userUsecase *usecase.UserUsecase
}

func NewUserController(userUsecase *usecase.UserUsecase) *UserController {
	return &UserController{userUsecase}
}

// GetProfile godoc
// @Summary Get User Profile
// @Description Fetch logged-in user profile details including points
// @Tags User Profile
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.BaseResponse{data=dto.UserProfileResponse}
// @Failure 401 {object} dto.BaseResponse
// @Router /api/v1/profile [get]
func (c *UserController) GetProfile(ctx *gin.Context) {
	// Ambil user_id dari middleware JWT
	userID := ctx.GetString("user_id")

	res, err := c.userUsecase.GetProfile(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.ErrorResponse(http.StatusNotFound, err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse(http.StatusOK, "Profile retrieved successfully", res))
}

// UpdateProfile godoc
// @Summary Update User Profile
// @Description Update name and username of the logged-in user
// @Tags User Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdateProfileRequest true "Profile Details"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseResponse
// @Router /api/v1/profile [put]
func (c *UserController) UpdateProfile(ctx *gin.Context) {
	userID := ctx.GetString("user_id")

	var req dto.UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(http.StatusBadRequest, "Invalid request body", err.Error()))
		return
	}

	err := c.userUsecase.UpdateProfile(ctx.Request.Context(), userID, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse(http.StatusOK, "Profile updated successfully", nil))
}