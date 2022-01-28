package handlers

import (
	"github.com/MiG-21/go-sso/internal"
	"github.com/MiG-21/go-sso/internal/sso"
	"github.com/MiG-21/go-sso/internal/web/types"
	"github.com/gofiber/fiber/v2"
)

// UserInfoHandler godoc
// @Summary user info
// @Description user info
// @Id user-info
// @Tags userinfo
// @Param token header string true "bearer token"
// @Accept json
// @Produce json
// @Success 200 {object} types.UserInfoResponse
// @Failure 400 {object} fiber.Error
// @Failure 422 {object} fiber.Error
// @Router /user/me [get]
func UserInfoHandler(s sso.SSOer, validator *internal.ServiceValidator) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		ctx.Set("Content-Type", "application/json")
		out := types.UserInfoResponse{}
		return ctx.Status(fiber.StatusOK).JSON(out)
	}
}

// RegisterHandler godoc
// @Summary register user
// @Description register user
// @Id register-user
// @Tags register
// @Param user body types.UserCreateRequest true "request body"
// @Accept json
// @Produce json
// @Success 200 {object} types.UserCreateResponse
// @Failure 400 {object} fiber.Error
// @Failure 422 {object} fiber.Error
// @Router /user/register [post]
func RegisterHandler(s sso.SSOer, validator *internal.ServiceValidator) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		params := &types.UserCreateRequest{}
		if err := ctx.BodyParser(params); err != nil {
			return HttpError(ctx, fiber.StatusBadRequest, err)
		}

		errors := HandleValidation(validator.Validate(params))
		if errors != nil {
			return HttpError(ctx, fiber.StatusUnprocessableEntity, errors)
		}
		out := types.UserCreateResponse{}
		return ctx.Status(fiber.StatusOK).JSON(out)
	}
}
