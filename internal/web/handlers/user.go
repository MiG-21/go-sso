package handlers

import (
	"time"

	"github.com/MiG-21/go-sso/internal"
	"github.com/MiG-21/go-sso/internal/event"
	"github.com/MiG-21/go-sso/internal/sso"
	"github.com/MiG-21/go-sso/internal/web/types"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
// @Failure 404 {object} fiber.Error
// @Failure 422 {object} fiber.Error
// @Failure 500 {object} fiber.Error
// @Router /user/me [post]
func UserInfoHandler(s sso.SSOer) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		ctx.Set("Content-Type", "application/json")
		claims := CtxClaims(ctx)
		user, err := s.UserManager().ById(claims.Id)
		if err != nil {
			return HttpError(ctx, fiber.StatusInternalServerError, err)
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
// @Success 201 {object} types.UserCreateResponse
// @Failure 422 {object} fiber.Error
// @Failure 500 {object} fiber.Error
// @Router /user/register [post]
func RegisterHandler(config *internal.Config, s sso.SSOer, validator *internal.ServiceValidator, eventService *event.Service) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		params := &types.UserCreateRequest{}
		if err := ctx.BodyParser(params); err != nil {
			return HttpError(ctx, fiber.StatusBadRequest, err)
		}

		errors := HandleValidation(validator.Validate(params))
		if errors != nil {
			return HttpError(ctx, fiber.StatusUnprocessableEntity, errors)
		}
		rand, err := uuid.NewRandom()
		if err != nil {
			return HttpError(ctx, fiber.StatusInternalServerError, errors)
		}
		user := &sso.UserModel{
			Name:     params.Name,
			Email:    params.Email,
			Password: internal.GetPasswordHash([]byte(params.Password)),
			Gender:   params.Gender,
			Active:   false,
			Locked:   false,
			Created:  time.Now().Unix(),
			Code:     rand.String(),
		}
		if err = s.UserManager().Validate(user); err != nil {
			return HttpError(ctx, fiber.StatusUnprocessableEntity, err)
		}
		if err = s.UserManager().Create(user); err != nil {
			return HttpError(ctx, fiber.StatusInternalServerError, err)
		}

		vUrl, err := user.GetVerificationUrl(ctx, config.Crypto.PrivateKey)
		if err != nil {
			return HttpError(ctx, fiber.StatusInternalServerError, err)
		}

		// emit event
		eventService.Emit(&event.UserCreated{
			UserName:        user.Name,
			UserEmail:       user.Email,
			VerificationUrl: vUrl,
		})

		out := types.UserCreateResponse{
			Name:  user.Name,
			Email: user.Email,
		}
		return ctx.Status(fiber.StatusCreated).JSON(out)
	}
}

// VerificationHandler godoc
// @Summary verify user
// @Description verify user
// @Id verify-user
// @Tags user
// @Param token query types.UserVerificationRequest true "request query"
// @Accept json
// @Produce json
// @Success 302 {string} string "Done"
// @Failure 422 {object} fiber.Error
// @Failure 500 {object} fiber.Error
// @Router /user/verification [get]
func VerificationHandler(config *internal.Config, s sso.SSOer, validator *internal.ServiceValidator, eventService *event.Service) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		params := &types.UserVerificationRequest{}
		if err := ctx.QueryParser(params); err != nil {
			return HttpError(ctx, fiber.StatusBadRequest, err)
		}

		errors := HandleValidation(validator.Validate(params))
		if errors != nil {
			return HttpError(ctx, fiber.StatusUnprocessableEntity, errors)
		}

		parsedToken, err := jwt.ParseWithClaims(params.Token, &internal.VerificationClaims{}, func(token *jwt.Token) (interface{}, error) {
			// since we only use the one private key to sign the tokens,
			// we also only use its public counterpart to verify
			return config.Crypto.PublicKey, nil
		})
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		claims, ok := parsedToken.Claims.(*internal.VerificationClaims)
		if !ok || !parsedToken.Valid {
			return fiber.NewError(fiber.StatusBadRequest, "invalid verification token")
		}

		ok, err = s.UserManager().Verify(claims.Id)
		if err != nil {
			return HttpError(ctx, fiber.StatusInternalServerError, err)
		}
		if !ok {
			return fiber.NewError(fiber.StatusBadRequest, "user not found")
		}

		return ctx.Redirect("Done", fiber.StatusFound)
	}
}
