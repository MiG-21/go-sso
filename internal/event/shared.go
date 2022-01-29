package event

import "github.com/rs/zerolog"

func SetupEventService(logger *zerolog.Logger) *Service {
	return &Service{logger: logger}
}
