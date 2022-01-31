package handlers

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type (
	ValidationError struct {
		Field string
		Tag   string
		Value string
	}
)

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s %s", e.Field, e.Tag)
}

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

func HandleValidation(errs error) []*ValidationError {
	var errors []*ValidationError
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			var element ValidationError
			element.Field = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

func ValidationErrorsToErrors(validationErrors []*ValidationError) []error {
	var errs []error
	for _, e := range validationErrors {
		errs = append(errs, e)
	}
	return errs
}
