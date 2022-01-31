package mail

import (
	"errors"
	"html/template"

	"github.com/MiG-21/go-sso/internal/event"
)

type (
	Service struct {
		EmailFrom string
		sender    EmailSender

		verificationEmailTpl    *template.Template
		passwordRecoverEmailTpl *template.Template
	}
)

func (s *Service) Sender() EmailSender {
	return s.sender
}

func (s *Service) SendActivationEmail(i interface{}) error {
	if e, ok := i.(*event.UserCreated); ok {
		data := struct {
			Name            string
			VerificationUrl string
		}{
			Name:            e.UserName,
			VerificationUrl: e.VerificationUrl.String(),
		}
		m := NewTemplate(s.EmailFrom, "Activation email", data, s.verificationEmailTpl, e.UserEmail)
		return s.sender.Send(m)
	}
	return errors.New("input data should be UserModel")
}

func (s *Service) SendPasswordRecoverEmail(i interface{}) error {
	if e, ok := i.(*event.UserPasswordRecover); ok {
		data := struct {
			Name            string
			VerificationUrl string
		}{
			Name:            e.UserName,
			VerificationUrl: e.VerificationUrl.String(),
		}
		m := NewTemplate(s.EmailFrom, "Password recover email", data, s.passwordRecoverEmailTpl, e.UserEmail)
		return s.sender.Send(m)
	}
	return errors.New("input data should be UserModel")
}
