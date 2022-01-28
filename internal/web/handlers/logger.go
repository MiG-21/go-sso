package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

const (
	ctxLoggerLocalsKey = "_ctx_locals_logger_"
)

func GetCtxLogger(ctx *fiber.Ctx) *zerolog.Logger {
	if l := ctx.Locals(ctxLoggerLocalsKey); l != nil {
		return l.(*zerolog.Logger)
	}
	return nil
}

func Logger(logger *zerolog.Logger) fiber.Handler {
	return func(ctx *fiber.Ctx) (err error) {
		startTime := time.Now()

		ctxLogger := logger.With().Str("requestId", ctx.Locals("requestid").(string)).Logger()

		// Add logger to context
		_ = ctx.Locals(ctxLoggerLocalsKey, &ctxLogger)

		// Handle request, store err for logging, manually call error handler
		if chainErr := ctx.Next(); chainErr != nil {
			if err = fiber.DefaultErrorHandler(ctx, chainErr); err != nil {
				_ = ctx.SendStatus(fiber.StatusInternalServerError)
			}
		}

		path := ctx.Context().Path()

		var env *zerolog.Event
		if ctx.Response().StatusCode() < 400 {
			env = ctxLogger.Info()
		} else {
			env = ctxLogger.Warn()
		}

		env.Int("code", ctx.Response().StatusCode()).
			Dur("time", time.Since(startTime)).
			Bytes("method", ctx.Context().Method()).
			Bytes("path", path).
			Str("addr", ctx.Context().RemoteAddr().String()).
			Send()

		return nil
	}
}
