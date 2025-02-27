package http

import "github.com/Dokito555/robin-notification/internal/services"

type NotificationController struct {
	NotificationService *services.NotificationService
}

func NewNotificationController(NotificationService *services.NotificationService) *NotificationController {
	return &NotificationController{
		NotificationService: NotificationService,
	}
}