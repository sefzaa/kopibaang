package controller

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"kopibang/domain/dto"
	"kopibang/usecase"
)

type VoucherController struct {
	voucherUsecase *usecase.VoucherUsecase
}

func NewVoucherController(voucherUsecase *usecase.VoucherUsecase) *VoucherController {
	return &VoucherController{voucherUsecase}
}

// CreateVoucher godoc
// @Summary Create Voucher
// @Description Admin can create a new discount voucher
// @Tags Admin Voucher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.VoucherRequest true "Voucher Details"
// @Success 201 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseResponse
// @Router /api/v1/admin/vouchers [post]
func (c *VoucherController) CreateVoucher(ctx *gin.Context) {
	var req dto.VoucherRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(http.StatusBadRequest, "Invalid request body", err.Error()))
		return
	}

	err := c.voucherUsecase.CreateVoucher(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusCreated, dto.SuccessResponse(http.StatusCreated, "Voucher created successfully", nil))
}

// UpdateVoucher godoc
// @Summary Update Voucher
// @Description Admin can update an existing voucher
// @Tags Admin Voucher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Voucher ID"
// @Param request body dto.VoucherRequest true "Voucher Details"
// @Success 200 {object} dto.BaseResponse
// @Router /api/v1/admin/vouchers/{id} [put]
func (c *VoucherController) UpdateVoucher(ctx *gin.Context) {
	id := ctx.Param("id")
	var req dto.VoucherRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(http.StatusBadRequest, "Invalid request body", err.Error()))
		return
	}

	err := c.voucherUsecase.UpdateVoucher(ctx.Request.Context(), id, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse(http.StatusOK, "Voucher updated successfully", nil))
}

// DeleteVoucher godoc
// @Summary Delete Voucher
// @Description Admin can delete a voucher
// @Tags Admin Voucher
// @Produce json
// @Security BearerAuth
// @Param id path string true "Voucher ID"
// @Success 200 {object} dto.BaseResponse
// @Router /api/v1/admin/vouchers/{id} [delete]
func (c *VoucherController) DeleteVoucher(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.voucherUsecase.DeleteVoucher(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse(http.StatusOK, "Voucher deleted successfully", nil))
}

// GetAllVouchers godoc
// @Summary Get All Vouchers
// @Description Fetch all vouchers for admin management
// @Tags Admin Voucher
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.BaseResponse{data=[]dto.VoucherResponse}
// @Router /api/v1/admin/vouchers [get]
func (c *VoucherController) GetAllVouchers(ctx *gin.Context) {
	vouchers, err := c.voucherUsecase.GetAllVouchers(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse(http.StatusOK, "Vouchers retrieved", vouchers))
}