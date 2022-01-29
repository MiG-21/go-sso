package mail

type (
	EmailSender interface {
		Send(mail Mailer) error
	}

	Mailer interface {
		From() string
		To() []string
		Subject() string
		Message() ([]byte, error)
	}
)
