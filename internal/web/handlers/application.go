package handlers

import (
	"time"

	"github.com/MiG-21/go-sso/internal"
	"github.com/MiG-21/go-sso/internal/models"
	"github.com/MiG-21/go-sso/internal/web/types"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// CreateApplicationHandler godoc
// @Summary create application
// @Description create application
// @Id create-application
// @Tags application
// @Param Authorization header string true "bearer token"
// @Param application body types.ApplicationCreateRequest true "request body"
// @Accept json
// @Produce json
// @Success 201 {object} types.ApplicationCreateResponse
// @Failure 422 {object} fiber.Error
// @Failure 500 {object} fiber.Error
// @Router /application/create [post]
func CreateApplicationHandler(s models.SSOer, validator *internal.ServiceValidator) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		params := &types.ApplicationCreateRequest{}
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
		application := &models.ApplicationModel{
			Application: params.Application,
			Domain:      params.Domain,
			RedirectUrl: params.RedirectUrl,
			Created:     time.Now().Unix(),
			Code:        rand.String(),
		}
		if err = s.ApplicationManager().Create(application); err != nil {
			return HttpError(ctx, fiber.StatusInternalServerError, err)
		}

		out := types.ApplicationCreateResponse{
			Application: application.Application,
			Domain:      application.Domain,
			RedirectUrl: application.RedirectUrl,
			Code:        application.Code,
		}
		return ctx.Status(fiber.StatusCreated).JSON(out)
	}
}
