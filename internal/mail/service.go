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

		verificationEmailTpl *template.Template
	}
)

func (s *Service) Sender() EmailSender {
	return s.sender
}

func (s *Service) SendVerificationEmail(i interface{}) error {
	if e, ok := i.(*event.UserCreated); ok {
		data := struct {
			Name            string
			VerificationUrl string
		}{
			Name:            e.UserName,
			VerificationUrl: e.VerificationUrl.String(),
		}
		m := NewTemplate(s.EmailFrom, "Verification email", data, s.verificationEmailTpl, e.UserEmail)
		return s.sender.Send(m)
	}
	return errors.New("input data should be UserModel")
}
