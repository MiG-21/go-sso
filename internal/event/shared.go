package event

import "github.com/rs/zerolog"

const (
	UserCreatedEvent     = "user_created"
	PasswordRecoverEvent = "password_recover_request"
)

func SetupEventService(logger *zerolog.Logger) *Service {
	return &Service{logger: logger}
}
