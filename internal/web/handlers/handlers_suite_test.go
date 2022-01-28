package handlers_test

import (
	"testing"

	"github.com/MiG-21/go-sso/internal"
	"github.com/MiG-21/go-sso/internal/web/handlers"
	"github.com/gofiber/fiber/v2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHandlers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Handlers Suite")
}

var (
	app *fiber.App
)

var _ = BeforeSuite(func() {
	app = fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	appConfig := &internal.Config{
		GitHash:   "SomeGitHash",
		GitBranch: "SomeGitBranch",
		GitUrl:    "SomeGitUrl",
		AppName:   "SomeApp",
		Version:   "SomeVersion",
	}

	app.Get("/v1/healthcheck/ping", handlers.HealthPingHandler)
	app.Get("/v1/healthcheck/info", handlers.HealthInfoHandler(appConfig))

	go func() {
		err := app.Listen(":5059")
		Expect(err).NotTo(HaveOccurred())
	}()
})

var _ = AfterSuite(func() {
	if app != nil {
		_ = app.Shutdown()
	}
})
