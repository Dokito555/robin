package email

import (
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

type Email struct {
	To      string
	Subject string
	Body    string
}

func (e *Email) SendEmail(config *viper.Viper, log *logrus.Logger) error {
	fromEmail := config.GetString("SMTP_AUTH_EMAIL")
	smtpHost := config.GetString("SMTP_HOST")
	smtpPass := config.GetString("SMTP_AUTH_PASS")
	smtpPortStr := config.GetString("SMTP_PORT")

	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		log.WithError(err).Error("invalid SMTP port configuration")
		return fmt.Errorf("invalid SMTP port: %w", err)
	}

	mailer := gomail.NewMessage()
	mailer.SetHeader("To", e.To)
	mailer.SetHeader("From", fromEmail)
	mailer.SetHeader("Subject", e.Subject)
	mailer.SetBody("text/html", e.Body)

	dialer := gomail.NewDialer(smtpHost, smtpPort, fromEmail, smtpPass)

	if err := dialer.DialAndSend(mailer); err != nil {
		log.WithError(err).Warn("failed to send email")
		return fmt.Errorf("failed to send email: %w", err)
	}
	
	return nil
}