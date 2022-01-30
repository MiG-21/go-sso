package internal

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"go.uber.org/dig"
)

type (
	SetupLoggerResult struct {
		dig.Out

		Logger *zerolog.Logger
		Error  error `group:"errors"`
	}
)

func SetupLogger(config *Config) SetupLoggerResult {
	sr := SetupLoggerResult{}
	if l, err := zerolog.ParseLevel(config.Logger.Level); err != nil {
		sr.Error = err
		return sr
	} else {
		zerolog.SetGlobalLevel(l)
	}

	zerolog.TimestampFieldName = "T"

	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	output.FormatLevel = func(i interface{}) string {
		return fmt.Sprintf("| %-6s|", i)
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
		return fmt.Sprintf("%s", i)
	}

	l := zerolog.New(output).With().Timestamp().Logger()
	sr.Logger = &l
	return sr
}
