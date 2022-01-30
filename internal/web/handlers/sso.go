package handlers

import (
	"errors"
	"strings"
	"time"

	"github.com/MiG-21/go-sso/internal"
	"github.com/MiG-21/go-sso/internal/models"
	"github.com/MiG-21/go-sso/internal/web/types"
	"github.com/MiG-21/go-sso/internal/web/views"
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
func AuthTokenHandler(s models.SSOer, validator *internal.ServiceValidator) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		ctx.Set("Content-Type", "application/json")
		params := &types.AuthRequest{}
		if err := ctx.BodyParser(params); err != nil {
			return HttpError(ctx, fiber.StatusBadRequest, err)
		}

		validationErrors := HandleValidation(validator.Validate(params))
		if validationErrors != nil {
			return HttpError(ctx, fiber.StatusUnprocessableEntity, validationErrors)
		}

		if app, err := s.ApplicationManager().ByCode(params.Code); err != nil || app == nil {
			return HttpError(ctx, fiber.StatusBadRequest, err)
		}

		item, err := s.UserManager().Authenticate(params.Email, params.Password)
		if err != nil {
			return HttpError(ctx, fiber.StatusBadRequest, err)
		}
		if item == nil {
			return HttpError(ctx, fiber.StatusUnauthorized, err)
		}
		exp := time.Now().Add(time.Hour * time.Duration(s.CTValidHours())).UTC()
		token, _ := s.BuildJWTToken(item.Id, strings.Split(item.Role, ","), exp)
		out := types.UserTokenResponse{Token: token}
		return ctx.Status(fiber.StatusOK).JSON(out)
	}
}

func AuthCookieHandler(s models.SSOer, validator *internal.ServiceValidator) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		params := &types.AuthRequest{}
		if err := ctx.BodyParser(params); err != nil {
			data := views.ErrorViewData(fiber.StatusBadRequest, err.Error())
			return ctx.Render("error", data, "layout")
		}

		validationErrors := HandleValidation(validator.Validate(params))
		if validationErrors != nil {
			var errs []error
			for _, e := range validationErrors {
				errs = append(errs, e)
			}
			data := views.LoginFormViewData(params.Code, errs...)
			return ctx.Render("login_form", data, "layout")
		}

		app, err := s.ApplicationManager().ByCode(params.Code)
		if err != nil {
			data := views.ErrorViewData(fiber.StatusBadRequest, err.Error())
			return ctx.Render("error", data, "layout")
		}
		item, err := s.UserManager().Authenticate(params.Email, params.Password)
		if err != nil {
			data := views.ErrorViewData(fiber.StatusUnauthorized, err.Error())
			return ctx.Render("error", data, "layout")
		}
		if item == nil {
			data := views.LoginFormViewData(params.Code, errors.New("email or password is incorrect"))
			return ctx.Render("login_form", data, "layout")
		}
		exp := time.Now().Add(time.Hour * time.Duration(s.CTValidHours())).UTC()
		token, _ := s.BuildJWTToken(item.Id, strings.Split(item.Role, ","), exp)
		cookie := s.BuildCookie(token, exp, app.Domain)
		ctx.Cookie(cookie)

		return ctx.Redirect(app.RedirectUrl, fiber.StatusFound)
	}
}

func LoginFormHandler(validator *internal.ServiceValidator) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		params := &types.LoginLogoutRequest{}
		if err := ctx.QueryParser(params); err != nil {
			data := views.ErrorViewData(fiber.StatusBadRequest, err.Error())
			return ctx.Render("error", data, "layout")
		}

		validationErrors := HandleValidation(validator.Validate(params))
		if validationErrors != nil {
			data := views.ErrorViewData(fiber.StatusUnprocessableEntity, validationErrors[0].Error())
			return ctx.Render("error", data, "layout")
		}

		data := views.LoginFormViewData(params.Code)
		return ctx.Render("login_form", data, "layout")
	}
}

func LogoutHandler(s models.SSOer, validator *internal.ServiceValidator) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		params := &types.LoginLogoutRequest{}
		if err := ctx.QueryParser(params); err != nil {
			data := views.ErrorViewData(fiber.StatusBadRequest, err.Error())
			return ctx.Render("error", data, "layout")
		}

		validationErrors := HandleValidation(validator.Validate(params))
		if validationErrors != nil {
			data := views.ErrorViewData(fiber.StatusUnprocessableEntity, validationErrors[0].Error())
			return ctx.Render("error", data, "layout")
		}
		app, err := s.ApplicationManager().ByCode(params.Code)
		if err != nil {
			data := views.ErrorViewData(fiber.StatusBadRequest, err.Error())
			return ctx.Render("error", data, "layout")
		}
		exp := time.Now().Add(time.Hour * time.Duration(-1))
		cookie := s.Logout(exp, app.Domain)
		ctx.Cookie(cookie)
		return ctx.Redirect("/login", fiber.StatusFound)
	}
}
