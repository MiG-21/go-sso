package mail

import (
	"go.uber.org/dig"
	"html/template"
	"strings"

	"github.com/MiG-21/go-sso/internal"
	"github.com/MiG-21/go-sso/internal/event"
)

type (
	SetupResult struct {
		dig.Out

		EmailService *Service
		Error        error `group:"errors"`
	}
)

func SetupService(config *internal.Config, eventService *event.Service) SetupResult {
	sr := SetupResult{}

	dir := strings.TrimRight(config.Frontend.Path, "/") + "/template/email/"
	layoutPath := dir + "layout.html"

	verificationEmailPath := dir + "verification_code.html"
	verificationEmailTpl, err := template.New("layout").ParseFiles(verificationEmailPath, layoutPath)
	if err != nil {
		sr.Error = err
		return sr
	}

	passwordRecoverEmailPath := dir + "password_recover.html"
	passwordRecoverEmailTpl, err := template.New("layout").ParseFiles(passwordRecoverEmailPath, layoutPath)
	if err != nil {
		sr.Error = err
		return sr
	}

	sender := &Smtp{
		host:     config.Smtp.SmtpHost,
		port:     config.Smtp.SmtpPort,
		user:     config.Smtp.SmtpUser,
		password: config.Smtp.SmtpPassword,
		useSSL:   config.Smtp.SmtpSsl,
	}

	service := &Service{
		EmailFrom:               config.Smtp.SmtpUser,
		sender:                  sender,
		verificationEmailTpl:    verificationEmailTpl,
		passwordRecoverEmailTpl: passwordRecoverEmailTpl,
	}

	eventService.AddListener(event.UserCreatedEvent, service.SendActivationEmail)
	eventService.AddListener(event.PasswordRecoverEvent, service.SendPasswordRecoverEmail)

	sr.EmailService = service

	return sr
}
