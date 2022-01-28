package internal

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"go.uber.org/dig"
)

type (
	AppRuntimeParams struct {
		dig.In

		App    *fiber.App
		Config *Config
		Logger *zerolog.Logger
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

func SetupLogger(config *Config) *zerolog.Logger {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if config.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	zerolog.TimestampFieldName = "T"

	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	output.FormatMessage = func(i interface{}) string {
		if i != nil {
			return fmt.Sprintf("***%s***", i)
		}
		return "--"
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%s", i))
	}

	l := zerolog.New(output).With().Timestamp().Logger()
	return &l
}

func Bootstrap(p AppRuntimeParams) {
	defer func() {
		// shutdown app
		if err := p.App.Shutdown(); err != nil {
			p.Logger.Err(err)
		}
	}()

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
			p.Logger.Err(err)
		}
	}()

	// handle shutdown
	quit := make(chan os.Signal, 2)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT) // subscribe to system quit
	<-quit
}
