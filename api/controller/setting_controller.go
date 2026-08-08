package controller

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"kopibang/domain/dto"
	"kopibang/usecase"
)

type SettingController struct {
	settingUsecase *usecase.SettingUsecase
}

func NewSettingController(settingUsecase *usecase.SettingUsecase) *SettingController {
	return &SettingController{settingUsecase}
}

// PatchBaristaStatus godoc
// @Summary Update Barista Status
// @Description Admin toggle to set available (online) or not available (offline)
// @Tags Admin Settings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.BaristaStatusRequest true "Availability Status"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseResponse
// @Router /api/v1/admin/settings/barista-status [patch]
func (c *SettingController) PatchBaristaStatus(ctx *gin.Context) {
	var req dto.BaristaStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(http.StatusBadRequest, "Invalid request body", err.Error()))
		return
	}

	err := c.settingUsecase.UpdateBaristaStatus(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse(http.StatusOK, "Barista status updated", nil))
}

// GetBaristaStatus godoc
// @Summary Get Barista Status
// @Description Customer/User checks if barista is currently available
// @Tags Settings
// @Produce json
// @Success 200 {object} dto.BaseResponse{data=dto.BaristaStatusResponse}
// @Router /api/v1/settings/barista-status [get]
func (c *SettingController) GetBaristaStatus(ctx *gin.Context) {
	res, _ := c.settingUsecase.GetBaristaStatus(ctx.Request.Context())
	ctx.JSON(http.StatusOK, dto.SuccessResponse(http.StatusOK, "Status retrieved", res))
}