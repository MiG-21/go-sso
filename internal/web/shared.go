package web

import (
	"time"

	"github.com/MiG-21/go-sso/internal"
	"github.com/MiG-21/go-sso/internal/sso"
	"github.com/MiG-21/go-sso/internal/web/handlers"
	swagger "github.com/arsmn/fiber-swagger/v2"
	goJson "github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/rs/zerolog"
	"go.uber.org/dig"

	// docs are generated by Swag CLI, you have to import them.
	// replace with your own docs folder, usually "github.com/username/reponame/docs"
	_ "github.com/MiG-21/go-sso/api/docs"
)

type (
	InitServerParams struct {
		dig.In

		Config    *internal.Config
		Logger    *zerolog.Logger
		Validator *internal.ServiceValidator
		Sso       sso.SSOer
	}
)

func SetupServer(p InitServerParams) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:         p.Config.AppName,
		ReadBufferSize:  p.Config.Http.ReadBufferSize, // Make sure these are big enough.
		WriteBufferSize: p.Config.Http.WriteBufferSize,
		ReadTimeout:     time.Duration(p.Config.Http.ReadTimeout) * time.Second,
		WriteTimeout:    time.Duration(p.Config.Http.WriteTimeout) * time.Second,
		IdleTimeout:     time.Duration(p.Config.Http.IdleTimeout) * time.Second, // This can be long for keep-alive connections.
		// DisableHeaderNormalizing:  true, // If you're not going to look at headers or know the casing you can set this.
		// DisableDefaultContentType: true, // Don't send Content-Type: text/plain if no Content-Type is set manually.
		DisableStartupMessage: true,
		JSONDecoder:           goJson.Unmarshal,
		JSONEncoder:           goJson.Marshal,
	})

	// Default middleware fiberApp
	app.Use(recover.New())
	// requestId middleware
	app.Use(requestid.New())
	// logger middleware
	app.Use(handlers.Logger(p.Logger))

	app.Static("/", p.Config.Frontend.Path, fiber.Static{
		Compress:      true,
		ByteRange:     true,
		Browse:        true,
		Index:         p.Config.Frontend.Index,
		CacheDuration: 10 * time.Second,
		MaxAge:        3600,
	})

	versionGroup := app.Group("v1")

	healthGroup := versionGroup.Group("healthcheck")
	healthGroup.Get("/ping", handlers.HealthPingHandler)
	healthGroup.Get("/info", handlers.HealthInfoHandler(p.Config))

	versionGroup.Post("/sso", handlers.AuthCookieHandler(p.Sso, p.Validator))
	versionGroup.Get("/logout", handlers.LogoutHandler(p.Sso))
	versionGroup.Post("/auth_token", handlers.AuthTokenHandler(p.Sso, p.Validator))

	userGroup := versionGroup.Group("user")
	userGroup.Post("/register", handlers.RegisterHandler(p.Sso, p.Validator))
	userGroup.Post("/lock", handlers.UserInfoHandler(p.Sso, p.Validator))
	userGroup.Post("/me", handlers.UserInfoHandler(p.Sso, p.Validator))
	userGroup.Post("/verify", handlers.UserInfoHandler(p.Sso, p.Validator))

	app.Get("/swagger/*", swagger.HandlerDefault) // default

	return app
}