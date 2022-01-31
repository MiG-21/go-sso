package internal

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/MiG-21/go-sso/internal/event"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"go.uber.org/dig"
)

type (
	AppRuntimeParams struct {
		dig.In

		App          *fiber.App
		Config       *Config
		Logger       *zerolog.Logger
		EventService *event.Service
		Errors       []error `group:"errors"`
	}

	ServiceValidator struct {
		Validator *validator.Validate
	}
)

func (sv *ServiceValidator) Validate(i interface{}) error {
	return sv.Validator.Struct(i)
}

func SetupValidator() *ServiceValidator {
	return &ServiceValidator{Validator: validator.New()}
}

func Bootstrap(p AppRuntimeParams) {
	terminate := false
	for _, err := range p.Errors {
		if err != nil {
			if p.Logger != nil {
				p.Logger.Err(err).Send()
			} else {
				log.Print(err)
			}
			terminate = true
		}
	}
	if terminate {
		p.Logger.Fatal().Msg("failed to start application")
	}

	defer func() {
		// shutdown events loop
		if err := p.EventService.Shutdown(); err != nil {
			p.Logger.Err(err).Send()
		}
		// shutdown app
		if err := p.App.Shutdown(); err != nil {
			p.Logger.Err(err).Send()
		}
	}()

	if err := p.EventService.Listen(); err != nil {
		p.Logger.Fatal().Err(err).Send()
	}

	go func() {
		config := p.App.Config()

		p.Logger.Info().
			Int("GOMAXPROCS", runtime.GOMAXPROCS(0)).
			Bool("isChild", fiber.IsChild()).
			Bool("prefork", config.Prefork).
			Int("pid", os.Getpid()).
			Str("appName", config.AppName).
			Int("port", p.Config.Port).
			Msg("Starting server...")

		if err := p.App.Listen(fmt.Sprintf(":%d", p.Config.Port)); err != nil {
			p.Logger.Fatal().Err(err).Send()
		}
	}()

	// handle shutdown
	quit := make(chan os.Signal, 2)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT) // subscribe to system quit
	<-quit
}
