package service

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/stonith404/pocket-id/backend/internal/common"
	"github.com/stonith404/pocket-id/backend/internal/utils/email"
	htemplate "html/template"
	"mime/multipart"
	"mime/quotedprintable"
	"net/smtp"
	"net/textproto"
	ttemplate "text/template"
)

type EmailService struct {
	appConfigService *AppConfigService
	htmlTemplates    map[string]*htemplate.Template
	textTemplates    map[string]*ttemplate.Template
}

func NewEmailService(appConfigService *AppConfigService) (*EmailService, error) {
	//TODO: -> config
	var templateDir = "./email-templates/"

	htmlTemplates, err := email.PrepareHTMLTemplates(templateDir, emailTemplatesPaths)
	if err != nil {
		return nil, fmt.Errorf("prepare html templates: %w", err)
	}

	textTemplates, err := email.PrepareTextTemplates(templateDir, emailTemplatesPaths)
	if err != nil {
		return nil, fmt.Errorf("prepare html templates: %w", err)
	}

	return &EmailService{
		appConfigService: appConfigService,
		htmlTemplates:    htmlTemplates,
		textTemplates:    textTemplates,
	}, nil
}

func SendEmail[V any](srv *EmailService, toEmail email.Address, template email.Template[V], tData *V) error {
	// Check if SMTP settings are set
	if srv.appConfigService.DbConfig.EmailEnabled.Value != "true" {
		return errors.New("email not enabled")
	}

	data := &email.TemplateData[V]{
		AppName: srv.appConfigService.DbConfig.AppName.Value,
		LogoURL: common.EnvConfig.AppURL + "/api/application-configuration/logo",
		Data:    tData,
	}

	body, boundary, err := prepareBody(srv, template, data)
	if err != nil {
		return fmt.Errorf("prepare email body for '%s': %w", template.Path, err)
	}

	// Construct the email message
	c := email.NewComposer()
	c.AddHeader("Subject", template.Title(data))
	c.AddAddressHeader("From", []email.Address{
		{
			Email: srv.appConfigService.DbConfig.SmtpFrom.Value,
			Name:  srv.appConfigService.DbConfig.AppName.Value,
		},
	})
	c.AddAddressHeader("To", []email.Address{toEmail})
	c.AddHeaderRaw("Content-Type",
		fmt.Sprintf("multipart/alternative;\n boundary=%s;\n charset=UTF-8", boundary),
	)
	c.Body(body)

	// Set up the authentication information.
	auth := smtp.PlainAuth("",
		srv.appConfigService.DbConfig.SmtpUser.Value,
		srv.appConfigService.DbConfig.SmtpPassword.Value,
		srv.appConfigService.DbConfig.SmtpHost.Value,
	)

	// Send the email
	err = smtp.SendMail(
		srv.appConfigService.DbConfig.SmtpHost.Value+":"+srv.appConfigService.DbConfig.SmtpPort.Value,
		auth,
		srv.appConfigService.DbConfig.SmtpFrom.Value,
		[]string{toEmail.Email},
		[]byte(c.String()),
	)

	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func prepareBody[V any](srv *EmailService, template email.Template[V], data *email.TemplateData[V]) (string, string, error) {
	body := bytes.NewBuffer(nil)
	mpart := multipart.NewWriter(body)

	// prepare text part
	var textHeader = textproto.MIMEHeader{}
	textHeader.Add("Content-Type", "text/plain;\n charset=UTF-8")
	textHeader.Add("Content-Transfer-Encoding", "quoted-printable")
	textPart, err := mpart.CreatePart(textHeader)
	if err != nil {
		return "", "", fmt.Errorf("create text part: %w", err)
	}

	textQp := quotedprintable.NewWriter(textPart)
	err = email.GetTemplate(srv.textTemplates, template).ExecuteTemplate(textQp, "root", data)
	if err != nil {
		return "", "", fmt.Errorf("execute text template: %w", err)
	}

	// prepare html part
	var htmlHeader = textproto.MIMEHeader{}
	htmlHeader.Add("Content-Type", "text/html;\n charset=UTF-8")
	htmlHeader.Add("Content-Transfer-Encoding", "quoted-printable")
	htmlPart, err := mpart.CreatePart(htmlHeader)
	if err != nil {
		return "", "", fmt.Errorf("create html part: %w", err)
	}

	htmlQp := quotedprintable.NewWriter(htmlPart)
	err = email.GetTemplate(srv.htmlTemplates, template).ExecuteTemplate(htmlQp, "root", data)
	if err != nil {
		return "", "", fmt.Errorf("execute html template: %w", err)
	}

	err = mpart.Close()
	if err != nil {
		return "", "", fmt.Errorf("close multipart: %w", err)
	}

	return body.String(), mpart.Boundary(), nil
}
