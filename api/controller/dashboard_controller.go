package controller

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"kopibang/domain/dto"
	"kopibang/usecase"
)

type DashboardController struct {
	dashboardUsecase *usecase.DashboardUsecase
}

func NewDashboardController(dashboardUsecase *usecase.DashboardUsecase) *DashboardController {
	return &DashboardController{dashboardUsecase}
}

// GetAdminDashboard godoc
// @Summary Admin Dashboard Metrics
// @Description Get dynamic metrics based on time filter (today, yesterday, this_week, this_month, this_year, custom)
// @Tags Admin Dashboard
// @Produce json
// @Security BearerAuth
// @Param filter query string false "Time Filter (default: today)" Enums(today, yesterday, this_week, this_month, this_year, custom)
// @Param start_date query string false "Start Date for custom filter (YYYY-MM-DD)"
// @Param end_date query string false "End Date for custom filter (YYYY-MM-DD)"
// @Success 200 {object} dto.BaseResponse{data=dto.DashboardResponse}
// @Failure 400 {object} dto.BaseResponse
// @Router /api/v1/admin/dashboard [get]
func (c *DashboardController) GetAdminDashboard(ctx *gin.Context) {
	var req dto.DashboardRequest
	// Binding parameter query ke struct DashboardRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(http.StatusBadRequest, "Invalid query parameters", err.Error()))
		return
	}

	res, err := c.dashboardUsecase.GetDashboardData(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse(http.StatusOK, "Dashboard data retrieved", res))
}