package controller

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"kopibang/domain/dto"
	"kopibang/internal/fcmutils"
	"firebase.google.com/go/v4/messaging"
)

type NotificationController struct {
	fcmClient *messaging.Client
}

func NewNotificationController(fcmClient *messaging.Client) *NotificationController {
	return &NotificationController{fcmClient}
}

// SendCustomPush godoc
// @Summary Send Custom Notification
// @Description Admin can send any message to all users
// @Tags Admin Notification
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CustomPushRequest true "Notification Data"
// @Success 200 {object} dto.BaseResponse
// @Router /api/v1/admin/notifications/push [post]
func (c *NotificationController) SendCustomPush(ctx *gin.Context) {
	var req dto.CustomPushRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(http.StatusBadRequest, "Invalid request body", err.Error()))
		return
	}

	fcmutils.SendToTopic(c.fcmClient, "all_users", req.Title, req.Body)

	ctx.JSON(http.StatusOK, dto.SuccessResponse(http.StatusOK, "Custom push notification sent successfully", nil))
}