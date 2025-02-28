package repository

import (
	"github.com/Dokito555/robin-notification/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type NotificationRepository struct {
	Repository[entity.NotificationHistory]
	Log *logrus.Logger
}


func NewNotificationRepository(log *logrus.Logger, db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{
		Repository: Repository[entity.NotificationHistory]{DB: db},
		Log: log,
	}
}

