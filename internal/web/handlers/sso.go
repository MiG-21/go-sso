package handlers

import (
	"time"

	"github.com/MiG-21/go-sso/internal"
	"github.com/MiG-21/go-sso/internal/sso"
	"github.com/MiG-21/go-sso/internal/web/types"
	"github.com/gofiber/fiber/v2"
)

// AuthTokenHandler godoc
// @Summary auth token
// @Description auth token
// @Id auth-token
// @Tags sso
// @Param params body types.AuthRequest true "request body"
// @Accept json
// @Produce json
// @Success 200 {object} types.UserTokenResponse
// @Failure 400 {object} fiber.Error
// @Failure 401 {object} fiber.Error
// @Failure 422 {object} fiber.Error
// @Router /auth_token [post]
func AuthTokenHandler(s sso.SSOer, validator *internal.ServiceValidator) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		ctx.Set("Content-Type", "application/json")
		params := &types.AuthRequest{}
		if err := ctx.BodyParser(params); err != nil {
			return HttpError(ctx, fiber.StatusBadRequest, err)
		}

		errors := HandleValidation(validator.Validate(params))
		if errors != nil {
			return HttpError(ctx, fiber.StatusUnprocessableEntity, errors)
		}
		item, err := s.UserManager().Authenticate(params.Email, params.Password)
		if err != nil {
			return HttpError(ctx, fiber.StatusBadRequest, err)
		}
		if item == nil {
			return HttpError(ctx, fiber.StatusUnauthorized, err)
		}
		exp := time.Now().Add(time.Hour * time.Duration(s.CTValidHours())).UTC()
		token, _ := s.BuildJWTToken(item.Id, nil, exp)
		out := types.UserTokenResponse{Token: token}
		return ctx.Status(fiber.StatusOK).JSON(out)
	}
}

// AuthCookieHandler godoc
// @Summary auth cookie
// @Description auth cookie
// @Id auth-cookie
// @Tags sso
// @Param params body types.AuthRequest true "request body"
// @Accept json
// @Produce json
// @Success 302 {string} string "Done"
// @Failure 400 {object} fiber.Error
// @Failure 422 {object} fiber.Error
// @Router /sso [post]
func AuthCookieHandler(s sso.SSOer, validator *internal.ServiceValidator) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		params := &types.AuthRequest{}
		if err := ctx.BodyParser(params); err != nil {
			return HttpError(ctx, fiber.StatusBadRequest, err)
		}

		errors := HandleValidation(validator.Validate(params))
		if errors != nil {
			return HttpError(ctx, fiber.StatusUnprocessableEntity, errors)
		}
		item, err := s.UserManager().Authenticate(params.Email, params.Password)
		if err != nil {
			return HttpError(ctx, fiber.StatusUnauthorized, err)
		}
		vh := s.CTValidHours()
		exp := time.Now().Add(time.Hour * time.Duration(vh)).UTC()
		token, _ := s.BuildJWTToken(item.Id, nil, exp)
		cookie := s.BuildCookie(token, exp)
		ctx.Cookie(cookie)

		return ctx.Redirect("Done", fiber.StatusFound)
	}
}

// LogoutHandler godoc
// @Summary logout user
// @Description logout user
// @Id logout-user
// @Tags sso
// @Accept json
// @Produce json
// @Success 302 {string} string "Done"
// @Router /logout [get]
func LogoutHandler(s sso.SSOer) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		exp := time.Now().Add(time.Hour * time.Duration(-1))
		cookie := s.Logout(exp)
		ctx.Cookie(cookie)
		return ctx.Redirect("Done", fiber.StatusFound)
	}
}
