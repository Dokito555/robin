package services

import (
	"bytes"
	"context"
	"html/template"

	"github.com/Dokito555/robin-notification/internal/entity"
	"github.com/Dokito555/robin-notification/internal/model"
	"github.com/Dokito555/robin-notification/internal/repository"
	"github.com/Dokito555/robin-notification/utils/errs"
	"github.com/sirupsen/logrus"
)

type NotificationService struct {
	Log                *logrus.Logger
	EmailRepository    *repository.EmailRepository
	TemplateRepository *repository.TemplateRepository
	MessagingService   *MessagingService
}

func NewNotificationService(log *logrus.Logger, emailRepo *repository.EmailRepository, templateRepo *repository.TemplateRepository, messaging *MessagingService) *NotificationService {
	return &NotificationService{
		Log:                log,
		EmailRepository:    emailRepo,
		TemplateRepository: templateRepo,
		MessagingService:   messaging,
	}
}

func (s *NotificationService) SendEmail(ctx context.Context, req *model.InternalNotificationRequest) error {
	s.Log.Info("starting Send Email function")
	s.Log.Infof("request received: %+v", req)

	emailmplt := new(entity.NotificationTemplate)
	err := s.TemplateRepository.GetTemplate(req.TemplateName, emailmplt)
	if err != nil {
		s.Log.Warn("failed to get template from database")
		return errs.ERROR_INTERNAL_SERVER_ERROR
	}

	tmpl, err := template.New("emailTemplate").Parse(emailmplt.Body)
	if err != nil {
		s.Log.Warn("failed to parse email template")
		return errs.ERROR_INTERNAL_SERVER_ERROR
	}

	var (
		tpl bytes.Buffer
	)

	err = tmpl.Execute(&tpl, req.Placeholder)
	if err != nil {
		s.Log.Warn("failed to execute placeholder")
		return errs.ERROR_INTERNAL_SERVER_ERROR
	}

	// email := proto.Email{
	// 	To:      req.Recipient,
	// 	Subject: emailTemplate.Subject,
	// 	Body:    tpl.String(),
	// }

	// err = email.SendEmail()
	// if err != nil {
	// 	notifHistory := &models.NotificationHistory{
	// 		Recipient:    req.Recipient,
	// 		TemplateID:   emailTemplate.ID,
	// 		Status:       "failed",
	// 		ErrorMessage: err.Error(),
	// 	}

	// 	s.EmailRepo.InsertNotificationHistory(ctx, notifHistory)
	// 	return errors.Wrap(err, "failed to send email")
	// }

	// notifHistory := &models.NotificationHistory{
	// 	Recipient:  req.Recipient,
	// 	TemplateID: emailTemplate.ID,
	// 	Status:     "SUCCESS",
	// }

	// s.EmailRepo.InsertNotificationHistory(ctx, notifHistory)

	// return nil
	return nil
}
