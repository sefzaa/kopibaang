package controller

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"kopibang/domain/dto"
	"kopibang/usecase"
)

type TransactionController struct {
	txUsecase *usecase.TransactionUsecase
}

func NewTransactionController(txUsecase *usecase.TransactionUsecase) *TransactionController {
	return &TransactionController{txUsecase}
}

// CreateOrder godoc
// @Summary Admin Input Order
// @Description Creates order, process redeem if any, returns receipt text and Earn QR
// @Tags Admin Transaction
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateOrderRequest true "Order Details"
// @Success 201 {object} dto.BaseResponse{data=dto.CreateOrderResponse}
// @Router /api/v1/admin/orders [post]
func (c *TransactionController) CreateOrder(ctx *gin.Context) {
	var req dto.CreateOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(http.StatusBadRequest, "Invalid request body", err.Error()))
		return
	}

	res, err := c.txUsecase.CreateOrder(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusCreated, dto.SuccessResponse(http.StatusCreated, "Order created successfully", res))
}

// RequestRedeemQR godoc
// @Summary User Request Redeem QR
// @Description User generates QR token to show to Admin
// @Tags User Transaction
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.BaseResponse{data=dto.RedeemQRResponse}
// @Router /api/v1/points/redeem-qr [post]
func (c *TransactionController) RequestRedeemQR(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	res, err := c.txUsecase.GenerateRedeemQR(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}
	ctx.JSON(http.StatusOK, dto.SuccessResponse(http.StatusOK, "QR generated", res))
}

// ScanEarnQR godoc
// @Summary User Scans Admin's QR
// @Description User scans QR generated from POS to earn points
// @Tags User Transaction
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.ScanEarnPointRequest true "Earn Token"
// @Success 200 {object} dto.BaseResponse
// @Router /api/v1/points/scan-earn [post]
func (c *TransactionController) ScanEarnQR(ctx *gin.Context) {
	var req dto.ScanEarnPointRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(http.StatusBadRequest, "Invalid request body", err.Error()))
		return
	}

	userID := ctx.GetString("user_id")
	err := c.txUsecase.ScanEarnQR(ctx.Request.Context(), userID, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse(http.StatusOK, "Points added successfully!", nil))
}