package handlers

import (
	"time"

	"github.com/MiG-21/go-sso/internal"
	"github.com/MiG-21/go-sso/internal/sso"
	"github.com/MiG-21/go-sso/internal/web/types"
	"github.com/gofiber/fiber/v2"
)

func AuthTokenHandler(s sso.SSOer, validator *internal.ServiceValidator) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		ctx.Set("Content-Type", "application/json")
		params := &types.AuthRequest{}
		if err := ctx.QueryParser(params); err != nil {
			return HttpError(ctx, fiber.StatusBadRequest, err)
		}

		errors := HandleValidation(validator.Validate(params))
		if errors != nil {
			return HttpError(ctx, fiber.StatusUnprocessableEntity, errors)
		}
		item, err := s.Login(params.Email, params.Password)
		if err != nil {
			return HttpError(ctx, fiber.StatusUnauthorized, err)
		}
		exp := time.Now().Add(time.Hour * time.Duration(s.CTValidHours())).UTC()
		token, _ := s.BuildJWTToken(item.Id, nil, exp)
		out := map[string]string{
			"token": token,
		}
		return ctx.Status(fiber.StatusOK).JSON(out)
	}
}

func AuthCookieHandler(s sso.SSOer, validator *internal.ServiceValidator) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		params := &types.AuthRequest{}
		if err := ctx.QueryParser(params); err != nil {
			return HttpError(ctx, fiber.StatusBadRequest, err)
		}

		errors := HandleValidation(validator.Validate(params))
		if errors != nil {
			return HttpError(ctx, fiber.StatusUnprocessableEntity, errors)
		}
		item, err := s.Login(params.Email, params.Password)
		if err != nil {
			return HttpError(ctx, fiber.StatusUnauthorized, err)
		}
		vh := s.CTValidHours()
		exp := time.Now().Add(time.Hour * time.Duration(vh)).UTC()
		token, _ := s.BuildJWTToken(item.Id, nil, exp)
		cookie := s.BuildCookie(token, exp)
		ctx.Cookie(cookie)

		return ctx.Redirect("", fiber.StatusFound)
	}
}

func LogoutHandler(s sso.SSOer) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		exp := time.Now().Add(time.Hour * time.Duration(-1))
		cookie := s.Logout(exp)
		ctx.Cookie(cookie)
		return ctx.Redirect("", fiber.StatusFound)
	}
}
