package handlers

import (
	"time"

	"github.com/MiG-21/go-sso/internal"
	"github.com/MiG-21/go-sso/internal/sso"
	"github.com/MiG-21/go-sso/internal/web/types"
	"github.com/gofiber/fiber/v2"
)

// UserInfoHandler godoc
// @Summary user info
// @Description user info
// @Id user-info
// @Tags user
// @Param Authorization header string true "bearer token"
// @Accept json
// @Produce json
// @Success 200 {object} types.UserInfoResponse
// @Failure 400 {object} fiber.Error
// @Failure 422 {object} fiber.Error
// @Router /user/me [post]
func UserInfoHandler(s sso.SSOer) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		ctx.Set("Content-Type", "application/json")
		claims := CtxClaims(ctx)
		user, err := s.UserManager().ById(claims.Id)
		if err != nil {
			return HttpError(ctx, fiber.StatusBadRequest, err)
		}
		if user == nil {
			return fiber.NewError(fiber.StatusNotFound)
		}
		out := types.UserInfoResponse{
			Id:       user.Id,
			Name:     user.Name,
			Email:    user.Email,
			Gender:   user.Gender,
			Data:     user.Data,
			Created:  user.Created,
			Updated:  user.Updated,
			Active:   user.Active,
			Locked:   user.Locked,
			LockedTo: user.LockedTo,
		}
		return ctx.Status(fiber.StatusOK).JSON(out)
	}
}

// RegisterHandler godoc
// @Summary register user
// @Description register user
// @Id register-user
// @Tags user
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
		user := &sso.UserModel{
			Name:     params.Name,
			Email:    params.Email,
			Password: internal.GetPasswordHash([]byte(params.Password)),
			Gender:   params.Gender,
			Active:   false,
			Locked:   false,
			Created:  time.Now().Unix(),
		}
		if err := s.UserManager().Validate(user); err != nil {
			return HttpError(ctx, fiber.StatusUnprocessableEntity, err)
		}
		if err := s.UserManager().Create(user); err != nil {
			return HttpError(ctx, fiber.StatusBadRequest, err)
		}
		out := types.UserCreateResponse{}
		return ctx.Status(fiber.StatusOK).JSON(out)
	}
}
