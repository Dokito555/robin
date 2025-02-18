package converter

import (
	"github.com/Dokito555/robin-notification/internal/entity"
	"github.com/Dokito555/robin-notification/internal/model"
)

func NotificationHistoryToResponse(n *entity.NotificationHistory) *model.NotificationHistory {
	return &model.NotificationHistory{
		ID:           n.ID,
		Recipient:    n.Recipient,
		TemplateID:   n.TemplateID,
		Status:       n.Status,
		ErrorMessage: n.ErrorMessage,
		CreatedAt:    n.CreatedAt,
		UpdatedAt:    n.UpdatedAt,
	}
}
