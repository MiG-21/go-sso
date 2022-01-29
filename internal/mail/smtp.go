package mail

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
)

type Smtp struct {
	host     string
	port     int
	user     string
	password string
	useSSL   bool
}

func (m *Smtp) Send(mail Mailer) error {
	if m.useSSL {
		return m.ssl(mail)
	} else {
		return m.def(mail)
	}
}

func (m *Smtp) ssl(mail Mailer) error {
	// TLS config
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         m.host,
	}

	// Connection
	conn, err := tls.Dial("tcp", m.server(), tlsConfig)
	if err != nil {
		return err
	}

	// Client
	client, err := smtp.NewClient(conn, m.host)
	if err != nil {
		return err
	}

	// Auth
	if err = client.Auth(m.auth()); err != nil {
		return err
	}

	// To && From
	if err = client.Mail(mail.From()); err != nil {
		return err
	}

	if err = client.Rcpt(mail.To()[0]); err != nil {
		return err
	}

	// Data
	w, err := client.Data()
	if err != nil {
		return err
	}

	// Body
	body, err := mail.Message()
	if err != nil {
		return err
	}

	_, err = w.Write(body)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return client.Quit()
}

func (m *Smtp) def(mail Mailer) error {
	body, err := mail.Message()
	if err != nil {
		return err
	}
	if err = smtp.SendMail(m.server(), m.auth(), mail.From(), mail.To(), body); err != nil {
		return err
	}
	return nil
}

func (m *Smtp) auth() smtp.Auth {
	// @TODO could be better implemented
	auth := smtp.PlainAuth("", m.user, m.password, m.host)
	return auth
}

func (m *Smtp) server() string {
	return m.host + ":" + fmt.Sprint(m.port)
}
