package views

import (
	"github.com/gofiber/fiber/v2"
)

func LoginFormViewData(code string, errs ...error) fiber.Map {
	return fiber.Map{
		"Code":   code,
		"Errors": errs,
	}
}

func PasswordRecoverFormViewData(code string, errs ...error) fiber.Map {
	return LoginFormViewData(code, errs...)
}

func ErrorViewData(code int, message string) fiber.Map {
	return fiber.Map{
		"Code":    code,
		"Message": message,
	}
}
