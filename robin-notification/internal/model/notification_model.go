package model

import "time"

type NotificationTemplate struct {
	ID           int       `json:"id"`
	TemplateName string    `json:"template_name"`
	Subject      string    `json:"subject""`
	Body         string    `json:"body"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type NotificationHistory struct {
	ID           int       `json:"id"`
	Recipient    string    `json:"recipient"`
	TemplateID   int       `json:"template_id"`
	Status       string    `json:"status"`
	ErrorMessage string    `json:"error_message"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type InternalNotificationRequest struct {
	TemplateName string `valid:"required"`
	Recipient    string `valid:"required"`
	Placeholder  map[string]string
}
