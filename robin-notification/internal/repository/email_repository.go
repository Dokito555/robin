package repository

import (
	"github.com/Dokito555/robin-notification/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type EmailRepository struct {
	Repository[entity.NotificationHistory]
	Log *logrus.Logger
}


func NewEmailRepository(log *logrus.Logger, db *gorm.DB) *EmailRepository {
	return &EmailRepository{
		Repository: Repository[entity.NotificationHistory]{DB: db},
		Log: log,
	}
}

