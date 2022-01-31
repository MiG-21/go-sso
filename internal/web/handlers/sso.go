package handlers

import (
	"errors"
	"strings"
	"time"

	"github.com/MiG-21/go-sso/internal"
	"github.com/MiG-21/go-sso/internal/event"
	"github.com/MiG-21/go-sso/internal/models"
	"github.com/MiG-21/go-sso/internal/web/types"
	"github.com/MiG-21/go-sso/internal/web/views"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
			data := views.LoginFormViewData(params.Code, ValidationErrorsToErrors(validationErrors)...)
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

func VerificationHandler(config *internal.Config, s models.SSOer, validator *internal.ServiceValidator) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		params := &types.UserVerificationRequest{}
		if err := ctx.QueryParser(params); err != nil {
			data := views.ErrorViewData(fiber.StatusBadRequest, err.Error())
			return ctx.Render("error", data, "layout")
		}

		validationErrors := HandleValidation(validator.Validate(params))
		if validationErrors != nil {
			data := views.ErrorViewData(fiber.StatusUnprocessableEntity, validationErrors[0].Error())
			return ctx.Render("error", data, "layout")
		}

		parsedToken, err := jwt.ParseWithClaims(params.Token, &internal.VerificationClaims{}, func(token *jwt.Token) (interface{}, error) {
			// since we only use the one private key to sign the tokens,
			// we also only use its public counterpart to verify
			return config.Crypto.PublicKey, nil
		})
		if err != nil {
			data := views.ErrorViewData(fiber.StatusBadRequest, err.Error())
			return ctx.Render("error", data, "layout")
		}

		claims, ok := parsedToken.Claims.(*internal.VerificationClaims)
		if !ok || !parsedToken.Valid || claims.Action != models.UserActionActivation {
			data := views.ErrorViewData(fiber.StatusBadRequest, "invalid verification token")
			return ctx.Render("error", data, "layout")
		}

		user, err := s.UserManager().ByCode(claims.Id)
		if err != nil {
			data := views.ErrorViewData(fiber.StatusBadRequest, err.Error())
			return ctx.Render("error", data, "layout")
		}
		if user == nil {
			data := views.ErrorViewData(fiber.StatusBadRequest, "invalid verification token")
			return ctx.Render("error", data, "layout")
		}

		user.Code = ""
		user.Active = true
		rows, err := s.UserManager().Update(user)
		if err != nil || rows == 0 {
			data := views.ErrorViewData(fiber.StatusBadRequest, "failed to activate user")
			return ctx.Render("error", data, "layout")
		}

		return ctx.Redirect("/verified", fiber.StatusFound)
	}
}

func VerifiedHandler() func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		data := struct {
			Message string
		}{
			Message: "Account has been verified",
		}
		return ctx.Render("landing", data, "layout")
	}
}

func PasswordRecoverFormHandler() func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		return ctx.Render("password_recover_form", nil, "layout")
	}
}

func PasswordRecoverHandler(config *internal.Config, s models.SSOer, validator *internal.ServiceValidator, eventService *event.Service) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		params := &types.PasswordRecoverRequest{}
		if err := ctx.BodyParser(params); err != nil {
			data := views.ErrorViewData(fiber.StatusBadRequest, err.Error())
			return ctx.Render("error", data, "layout")
		}

		validationErrors := HandleValidation(validator.Validate(params))
		if validationErrors != nil {
			data := views.ErrorViewData(fiber.StatusUnprocessableEntity, validationErrors[0].Error())
			return ctx.Render("error", data, "layout")
		}
		user, err := s.UserManager().ByEmail(params.Email)
		if err != nil {
			data := views.ErrorViewData(fiber.StatusBadRequest, err.Error())
			return ctx.Render("error", data, "layout")
		}
		if user == nil {
			data := views.ErrorViewData(fiber.StatusBadRequest, "invalid verification token")
			return ctx.Render("error", data, "layout")
		}

		rand, err := uuid.NewRandom()
		if err != nil {
			data := views.ErrorViewData(fiber.StatusBadRequest, err.Error())
			return ctx.Render("error", data, "layout")
		}
		user.Code = rand.String()
		_, err = s.UserManager().Update(user)
		if err != nil {
			data := views.ErrorViewData(fiber.StatusBadRequest, err.Error())
			return ctx.Render("error", data, "layout")
		}

		vUrl, err := user.GetActionUrl(ctx, "/password/change", models.UserActionPasswordRecover, config.Crypto.PrivateKey)
		if err != nil {
			return HttpError(ctx, fiber.StatusInternalServerError, err)
		}

		// emit event
		eventService.Emit(&event.UserPasswordRecover{
			UserName:        user.Name,
			UserEmail:       user.Email,
			VerificationUrl: vUrl,
		})

		return ctx.Redirect("/password/recover/send", fiber.StatusFound)
	}
}

func PasswordRecoverSendHandler() func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		data := struct {
			Message string
		}{
			Message: "Password recover link has been send, please, check your mailbox",
		}
		return ctx.Render("landing", data, "layout")
	}
}

func PasswordChangeFormHandler(config *internal.Config, validator *internal.ServiceValidator) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		params := &types.UserVerificationRequest{}
		if err := ctx.QueryParser(params); err != nil {
			data := views.ErrorViewData(fiber.StatusBadRequest, err.Error())
			return ctx.Render("error", data, "layout")
		}

		validationErrors := HandleValidation(validator.Validate(params))
		if validationErrors != nil {
			data := views.ErrorViewData(fiber.StatusUnprocessableEntity, validationErrors[0].Error())
			return ctx.Render("error", data, "layout")
		}

		parsedToken, err := jwt.ParseWithClaims(params.Token, &internal.VerificationClaims{}, func(token *jwt.Token) (interface{}, error) {
			// since we only use the one private key to sign the tokens,
			// we also only use its public counterpart to verify
			return config.Crypto.PublicKey, nil
		})
		if err != nil {
			data := views.ErrorViewData(fiber.StatusBadRequest, err.Error())
			return ctx.Render("error", data, "layout")
		}

		claims, ok := parsedToken.Claims.(*internal.VerificationClaims)
		if !ok || !parsedToken.Valid || claims.Action != models.UserActionPasswordRecover {
			data := views.ErrorViewData(fiber.StatusBadRequest, "invalid verification token")
			return ctx.Render("error", data, "layout")
		}
		data := views.PasswordRecoverFormViewData(claims.Id)
		return ctx.Render("password_change_form", data, "layout")
	}
}

func PasswordChangeHandler(s models.SSOer, validator *internal.ServiceValidator) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		params := &types.PasswordChangeRequest{}
		if err := ctx.BodyParser(params); err != nil {
			data := views.ErrorViewData(fiber.StatusBadRequest, err.Error())
			return ctx.Render("error", data, "layout")
		}

		validationErrors := HandleValidation(validator.Validate(params))
		if validationErrors != nil {
			data := views.PasswordRecoverFormViewData(params.Code, ValidationErrorsToErrors(validationErrors)...)
			return ctx.Render("password_change_form", data, "layout")
		}

		user, err := s.UserManager().ByCode(params.Code)
		if err != nil {
			data := views.ErrorViewData(fiber.StatusBadRequest, err.Error())
			return ctx.Render("error", data, "layout")
		}
		if user == nil {
			data := views.ErrorViewData(fiber.StatusBadRequest, "invalid verification token")
			return ctx.Render("error", data, "layout")
		}

		user.Code = ""
		user.Password = internal.GetPasswordHash([]byte(params.Password))
		rows, err := s.UserManager().Update(user)
		if err != nil || rows == 0 {
			data := views.ErrorViewData(fiber.StatusBadRequest, "failed to change password")
			return ctx.Render("error", data, "layout")
		}

		return ctx.Redirect("/login", fiber.StatusFound)
	}
}
