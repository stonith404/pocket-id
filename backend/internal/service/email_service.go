package service

import (
	"errors"
	"fmt"
	"github.com/stonith404/pocket-id/backend/internal/common"
	"net/smtp"
	"os"
	"strings"
)

type EmailService struct {
	appConfigService *AppConfigService
}

func NewEmailService(appConfigService *AppConfigService) *EmailService {
	return &EmailService{
		appConfigService: appConfigService}
}

// Send sends an email notification
func (s *EmailService) Send(toEmail, title, templateName string, templateParameters map[string]interface{}) error {
	// Check if SMTP settings are set
	if s.appConfigService.DbConfig.EmailEnabled.Value != "true" {
		return errors.New("email not enabled")
	}

	// Construct the email message
	subject := fmt.Sprintf("Subject: %s\n", title)
	subject += "From: " + s.appConfigService.DbConfig.SmtpFrom.Value + "\n"
	subject += "To: " + toEmail + "\n"
	subject += "Content-Type: text/html; charset=UTF-8\n"

	body, err := os.ReadFile(fmt.Sprintf("./email-templates/%s.html", templateName))
	bodyString := string(body)
	if err != nil {
		return fmt.Errorf("failed to read email template: %w", err)
	}

	// Replace template parameters
	templateParameters["appName"] = s.appConfigService.DbConfig.AppName.Value
	templateParameters["appUrl"] = common.EnvConfig.AppURL
	
	for key, value := range templateParameters {
		bodyString = strings.ReplaceAll(bodyString, fmt.Sprintf("{{%s}}", key), fmt.Sprintf("%v", value))
	}

	emailBody := []byte(subject + bodyString)

	// Set up the authentication information.
	auth := smtp.PlainAuth("", s.appConfigService.DbConfig.SmtpUser.Value, s.appConfigService.DbConfig.SmtpPassword.Value, s.appConfigService.DbConfig.SmtpHost.Value)

	// Send the email
	err = smtp.SendMail(
		s.appConfigService.DbConfig.SmtpHost.Value+":"+s.appConfigService.DbConfig.SmtpPort.Value,
		auth,
		s.appConfigService.DbConfig.SmtpFrom.Value,
		[]string{toEmail},
		emailBody,
	)

	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
