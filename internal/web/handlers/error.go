package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type (
	ErrorResponse struct {
		FailedField string
		Tag         string
		Value       string
	}
)

func HttpError(ctx *fiber.Ctx, status int, data interface{}) error {
	if ctx.Context() == nil {
		return nil
	}
	switch data.(type) {
	case error:
		return fiber.NewError(status, data.(error).Error())
	default:
		return ctx.Status(status).JSON(data)
	}
}

func HandleValidation(errs error) []*ErrorResponse {
	var errors []*ErrorResponse
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}
