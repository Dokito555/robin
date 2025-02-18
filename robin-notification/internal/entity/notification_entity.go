package entity

import "time"

type NotificationTemplate struct {
	ID           int       `gorm:"column:id;type:primaryKey"`
	TemplateName string    `gorm:"column:template_name"`
	Subject      string    `gorm:"column:subject"`
	Body         string    `gorm:"column:body"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

func (*NotificationTemplate) TableName() string {
	return "notification_templates"
}

type NotificationHistory struct {
	ID           int	`gorm:"column:id;type:primaryKey"`
	Recipient    string `gorm:"column:recipient"`
	TemplateID   int    `gorm:"column:template_id"`
	Status       string `gorm:"column:status"`
	ErrorMessage string `gorm:"column:error_message;type:text"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

func (*NotificationHistory) TableName() string {
	return "notification_history"
}