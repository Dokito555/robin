package repository

import (
	"github.com/Dokito555/robin-notification/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TemplateRepository struct {
	Repository[entity.NotificationTemplate]
	Log *logrus.Logger
}

func NewTemplateRepository(log *logrus.Logger, db *gorm.DB) *TemplateRepository {
	return &TemplateRepository{
		Repository: Repository[entity.NotificationTemplate]{DB: db},
		Log: log,
	}
}

func (r *TemplateRepository) GetTemplate(templateName string, template *entity.NotificationTemplate) (error) {
	return r.DB.Where("template_name = ?", templateName).First(template).Error
} 