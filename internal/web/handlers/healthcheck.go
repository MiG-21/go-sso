package handlers

import (
	"github.com/MiG-21/go-sso/internal"
	"github.com/MiG-21/go-sso/internal/web/types"
	"github.com/gofiber/fiber/v2"
)

// HealthPingHandler godoc
// @Summary bidBucket health checker ping
// @Description bidBucket health checker ping
// @Id health-check-ping
// @Tags healthcheck
// @Accept x-www-form-urlencoded
// @Produce json
// @Success 200 {object} types.HealthCheckPing
// @Router /healthcheck/ping [get]
func HealthPingHandler(ctx *fiber.Ctx) error {
	ctx.Set("Content-Type", "application/json")
	out := types.HealthCheckPing{
		Ping: "PONG",
	}
	return ctx.Status(fiber.StatusOK).JSON(out)
}

// HealthInfoHandler godoc
// @Summary bidBucket health checker info
// @Description bidBucket health checker info
// @Id health-check-info
// @Tags healthcheck
// @Accept x-www-form-urlencoded
// @Produce json
// @Success 200 {object} types.HealthCheckInfo
// @Router /healthcheck/info [get]
func HealthInfoHandler(cfg *internal.Config) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		ctx.Set("Content-Type", "application/json")
		out := types.HealthCheckInfo{
			AppName:    cfg.AppName,
			AppVersion: cfg.Version,
			Git: types.HealthCheckInfoGit{
				Hash: cfg.GitHash,
				Ref:  cfg.GitBranch,
				Url:  cfg.GitUrl,
			},
		}
		return ctx.Status(fiber.StatusOK).JSON(out)
	}
}
