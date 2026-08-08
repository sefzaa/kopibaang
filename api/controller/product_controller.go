package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"kopibang/domain/dto"
	"kopibang/usecase"
)

type ProductController struct {
	productUsecase *usecase.ProductUsecase
}

func NewProductController(productUsecase *usecase.ProductUsecase) *ProductController {
	return &ProductController{productUsecase}
}

// CreateMenu godoc
// @Summary Create New Coffee Menu
// @Description Admin can create a new menu with ingredients
// @Tags Admin Menu
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.ProductRequest true "Menu Details"
// @Success 201 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseResponse
// @Router /api/v1/admin/menus [post]
func (c *ProductController) CreateMenu(ctx *gin.Context) {
	var req dto.ProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(http.StatusBadRequest, "Invalid request body", err.Error()))
		return
	}

	err := c.productUsecase.CreateMenu(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusCreated, dto.SuccessResponse(http.StatusCreated, "Menu created successfully", nil))
}

// UpdateMenu godoc
// @Summary Update Existing Menu
// @Description Admin can update menu details and ingredients
// @Tags Admin Menu
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Param request body dto.ProductRequest true "Menu Details"
// @Success 200 {object} dto.BaseResponse
// @Router /api/v1/admin/menus/{id} [put]
func (c *ProductController) UpdateMenu(ctx *gin.Context) {
	id := ctx.Param("id")
	var req dto.ProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(http.StatusBadRequest, "Invalid request body", err.Error()))
		return
	}

	err := c.productUsecase.UpdateMenu(ctx.Request.Context(), id, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse(http.StatusOK, "Menu updated successfully", nil))
}

// DeleteMenu godoc
// @Summary Delete a Menu
// @Description Admin can delete a menu by ID
// @Tags Admin Menu
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Success 200 {object} dto.BaseResponse
// @Router /api/v1/admin/menus/{id} [delete]
func (c *ProductController) DeleteMenu(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.productUsecase.DeleteMenu(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse(http.StatusOK, "Menu deleted successfully", nil))
}

// GetMenus godoc
// @Summary Get All Menus
// @Description User can get all active menus, Admin gets all menus (active & inactive)
// @Tags Menu
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.BaseResponse{data=[]dto.ProductResponse}
// @Router /api/v1/menus [get]
func (c *ProductController) GetMenus(ctx *gin.Context) {
	role := ctx.GetString("role") // didapat dari JWT middleware

	menus, err := c.productUsecase.GetMenus(ctx.Request.Context(), role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse(http.StatusOK, "Success fetch menus", menus))
}

// GetMenuDetail godoc
// @Summary Get Menu Detail
// @Description Get specific menu by ID
// @Tags Menu
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Success 200 {object} dto.BaseResponse{data=dto.ProductResponse}
// @Router /api/v1/menus/{id} [get]
func (c *ProductController) GetMenuDetail(ctx *gin.Context) {
	id := ctx.Param("id")
	role := ctx.GetString("role")

	menu, err := c.productUsecase.GetMenuByID(ctx.Request.Context(), id, role)
	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.ErrorResponse(http.StatusNotFound, err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse(http.StatusOK, "Success fetch menu detail", menu))
}

// ToggleMenuStatus godoc
// @Summary Archive or Unarchive a Menu
// @Description Toggle the active status of a menu
// @Tags Menus
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseResponse
// @Router /api/v1/admin/menus/{id}/status [patch]
func (c *ProductController) ToggleMenuStatus(ctx *gin.Context) {
	productID := ctx.Param("id")

	err := c.productUsecase.ToggleMenuStatus(ctx.Request.Context(), productID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse(http.StatusOK, "Menu status toggled successfully", nil))
}

// UploadImage godoc
// @Summary Upload Menu Image to MinIO
// @Description Admin uploads an image, returns the public URL
// @Tags Admin Menu
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "Image File"
// @Success 200 {object} dto.BaseResponse
// @Router /api/v1/admin/upload [post]
func (c *ProductController) UploadImage(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(http.StatusBadRequest, "No file uploaded or invalid key (use 'file')", err.Error()))
		return
	}

	url, err := c.productUsecase.UploadImageToMinIO(ctx.Request.Context(), file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(http.StatusInternalServerError, "Failed to upload image", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse(http.StatusOK, "Image uploaded successfully", gin.H{"image_url": url}))
}