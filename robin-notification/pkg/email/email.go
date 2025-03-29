package email

import (
	"fmt"
	"strconv"

	"github.com/Dokito555/robin-notification/internal/entity"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

type EmailPkg struct {
	Log    *logrus.Logger
	Config *viper.Viper
}

func NewEmailPkg(log *logrus.Logger, config *viper.Viper) *EmailPkg {
	return &EmailPkg{
		Log:    log,
		Config: config,
	}
}

func (e *EmailPkg) SendEmail(mail *entity.Email) error {
	e.Log.Info("sending email through smtp server")
	e.Log.Info("email recieved: %v", mail)

	fromEmail := e.Config.GetString("SMTP_AUTH_EMAIL")
	smtpHost := e.Config.GetString("SMTP_HOST")
	smtpPass := e.Config.GetString("SMTP_AUTH_PASS")
	smtpPortStr := e.Config.GetString("SMTP_PORT")

	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		e.Log.WithError(err).Error("invalid SMTP port configuration")
		return fmt.Errorf("invalid SMTP port: %w", err)
	}

	mailer := gomail.NewMessage()
	mailer.SetHeader("To", mail.To)
	mailer.SetHeader("From", fromEmail)
	mailer.SetHeader("Subject", mail.Subject)
	mailer.SetBody("text/html", mail.Body)

	dialer := gomail.NewDialer(smtpHost, smtpPort, fromEmail, smtpPass)

	if err := dialer.DialAndSend(mailer); err != nil {
		e.Log.WithError(err).Warn("failed to send email")
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
