package services

import (
	// "bytes"
	"context"
	// "encoding/json"
	"fmt"
	// "html/template"
	// "time"

	"github.com/Dokito555/robin-notification/internal/entity"
	"github.com/Dokito555/robin-notification/internal/model"
	"github.com/Dokito555/robin-notification/internal/repository"
	"github.com/Dokito555/robin-notification/pkg/email"
	// "github.com/Dokito555/robin-notification/utils/errs"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type INotifier interface {
	SendEmail(ctx context.Context, req *model.InternalNotificationRequest) error
}

type NotificationService struct {
	DB                     *gorm.DB
	Log                    *logrus.Logger
	Config                 *viper.Viper
	EmailRepository        *repository.EmailRepository
	TemplateRepository     *repository.TemplateRepository
	NotificationRepository *repository.NotificationRepository
	EmailPkg               *email.EmailPkg
}

func NewNotificationService(db *gorm.DB, log *logrus.Logger, config *viper.Viper, emailRepo *repository.EmailRepository, notificationRepo *repository.NotificationRepository, templateRepo *repository.TemplateRepository, emailPkg *email.EmailPkg) *NotificationService {
	return &NotificationService{
		DB:                     db,
		Log:                    log,
		Config:                 config,
		EmailRepository:        emailRepo,
		TemplateRepository:     templateRepo,
		NotificationRepository: notificationRepo,
		EmailPkg:               emailPkg,
	}
}

// func (s *NotificationService) SendEmail(ctx context.Context, req *model.InternalNotificationRequest) error {
// 	s.Log.Info("starting Send Email function")
// 	s.Log.Infof("request received: %+v", req)

// 	emailmplt := new(entity.NotificationTemplate)
// 	err := s.TemplateRepository.GetTemplate(req.TemplateName, emailmplt)
// 	if err != nil {
// 		s.Log.Warn("failed to get template from database")
// 		return errs.ERROR_INTERNAL_SERVER_ERROR
// 	}

// 	msgChan, errChan, closeFunc, err := s.MessagingService.ConsumeKafkaMessage(s.Config.GetString("KAFKA_REGISTER_TOPIC"))
// 	if err != nil {
// 		s.Log.Warn("failed to consume kafka message")
// 		return errs.ERROR_INTERNAL_SERVER_ERROR
// 	}

// 	defer closeFunc()

// 	var registerMessage struct {
// 		Username string `json:"username"`
// 		Role     string `json:"role"`
// 	}

// 	select {
// 	case msg := <-msgChan:
// 		s.Log.Infof("kafka message revieved: %s: ", string(msg))
// 		if err := json.Unmarshal(msg, &registerMessage); err != nil {
// 			s.Log.Errorf("failed to unmarshal kafka message: %v", err)
// 			return errs.ERROR_INTERNAL_SERVER_ERROR
// 		}

// 		if req.Placeholder == nil {
// 			req.Placeholder = make(map[string]interface{})
// 		}
// 		req.Placeholder["username"] = registerMessage.Username
// 		req.Placeholder["role"] = registerMessage.Role
// 	case err := <-errChan:
// 		s.Log.Errorf("error from kafka consumer: %v", err)
// 		return errs.ERROR_INTERNAL_SERVER_ERROR
// 		// timeout 10 seconds on kafka consumer
// 	case <-time.After(10 * time.Second):
// 		s.Log.Warn("timeout waiting for kafka message")
// 		return errs.ERROR_INTERNAL_SERVER_ERROR
// 	case <-ctx.Done():
// 		s.Log.Warn("context cancelled before kafka message recieved")
// 		return errs.ERROR_INTERNAL_SERVER_ERROR
// 	}

// 	tmpl, err := template.New("emailTemplate").Parse(emailmplt.Body)
// 	if err != nil {
// 		s.Log.Warn("failed to parse email template")
// 		return errs.ERROR_INTERNAL_SERVER_ERROR
// 	}

// 	var (
// 		tpl bytes.Buffer
// 	)

// 	err = tmpl.Execute(&tpl, req.Placeholder)
// 	if err != nil {
// 		s.Log.Warn("failed to execute placeholder")
// 		return errs.ERROR_INTERNAL_SERVER_ERROR
// 	}

// 	newEmail := &entity.Email{
// 		To:      req.Recipient,
// 		Subject: emailmplt.Subject,
// 		Body:    tpl.String(),
// 	}

// 	err = s.EmailPkg.SendEmail(newEmail)
// 	if err != nil {
// 		s.Log.Warn("failed to send email via pkg")
// 		history := &entity.NotificationHistory{
// 			Recipient:    req.Recipient,
// 			TemplateID:   emailmplt.ID,
// 			Status:       "FAILED",
// 			ErrorMessage: err.Error(),
// 		}
// 		err = s.NotificationRepository.Create(s.DB, history)
// 		if err != nil {
// 			s.Log.Warn("failed to create failed email history in database")
// 			return errs.ERROR_INTERNAL_SERVER_ERROR
// 		}
// 		return errs.ERROR_INTERNAL_SERVER_ERROR
// 	}

// 	history := &entity.NotificationHistory{
// 		Recipient:    req.Recipient,
// 		TemplateID:   emailmplt.ID,
// 		Status:       "SUCCESS",
// 	}

// 	err = s.NotificationRepository.Create(s.DB, history)
// 	if err != nil {
// 		s.Log.Warn("failed to create email history in database")
// 		return errs.ERROR_INTERNAL_SERVER_ERROR
// 	}

// 	return nil
// }

func (s *NotificationService) SendEmail(ctx context.Context, req *model.InternalNotificationRequest) error {
	s.Log.Info("Starting SendEmail function")
	
	body := fmt.Sprintf("Hello %v, your role is %v", req.Placeholder["username"], req.Placeholder["role"])
	newEmail := &entity.Email{
		To:      req.Recipient,
		Subject: "Welcome Email",
		Body:    body,
	}

	if err := s.EmailPkg.SendEmail(newEmail); err != nil {
		s.Log.Errorf("Failed to send email: %v", err)
		return err
	}

	s.Log.Info("Email sent successfully")
	return nil
}