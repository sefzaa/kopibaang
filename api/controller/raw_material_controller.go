package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"kopibang/domain/dto"
	"kopibang/usecase"
)

type RawMaterialController struct {
	materialUsecase *usecase.RawMaterialUsecase
}

func NewRawMaterialController(materialUsecase *usecase.RawMaterialUsecase) *RawMaterialController {
	return &RawMaterialController{materialUsecase}
}

// AddMaterial godoc
// @Summary Add Raw Material
// @Description Admin can add new raw material stock
// @Tags Admin Raw Material
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.RawMaterialRequest true "Material Details"
// @Success 201 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseResponse
// @Router /api/v1/admin/raw-materials [post]
func (c *RawMaterialController) AddMaterial(ctx *gin.Context) {
	var req dto.RawMaterialRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(http.StatusBadRequest, "Invalid request body", err.Error()))
		return
	}

	err := c.materialUsecase.AddMaterial(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusCreated, dto.SuccessResponse(http.StatusCreated, "Raw material added successfully", nil))
}

// UpdateMaterial godoc
// @Summary Update Raw Material
// @Description Admin can update existing raw material details
// @Tags Admin Raw Material
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Material ID"
// @Param request body dto.RawMaterialRequest true "Material Details"
// @Success 200 {object} dto.BaseResponse
// @Router /api/v1/admin/raw-materials/{id} [put]
func (c *RawMaterialController) UpdateMaterial(ctx *gin.Context) {
	id := ctx.Param("id")
	var req dto.RawMaterialRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(http.StatusBadRequest, "Invalid request body", err.Error()))
		return
	}

	err := c.materialUsecase.UpdateMaterial(ctx.Request.Context(), id, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse(http.StatusOK, "Raw material updated successfully", nil))
}

// DeleteMaterial godoc
// @Summary Delete Raw Material
// @Description Admin can delete a raw material
// @Tags Admin Raw Material
// @Produce json
// @Security BearerAuth
// @Param id path string true "Material ID"
// @Success 200 {object} dto.BaseResponse
// @Router /api/v1/admin/raw-materials/{id} [delete]
func (c *RawMaterialController) DeleteMaterial(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.materialUsecase.DeleteMaterial(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse(http.StatusOK, "Raw material deleted successfully", nil))
}

// GetAllMaterials godoc
// @Summary Get All Raw Materials
// @Description Fetch all raw material stock lists
// @Tags Admin Raw Material
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.BaseResponse{data=[]dto.RawMaterialResponse}
// @Router /api/v1/admin/raw-materials [get]
func (c *RawMaterialController) GetAllMaterials(ctx *gin.Context) {
	materials, err := c.materialUsecase.GetAllMaterials(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse(http.StatusOK, "Raw materials retrieved", materials))
}