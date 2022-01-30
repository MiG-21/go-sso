package views

import (
	"github.com/gofiber/fiber/v2"
)

func LoginFormViewData(code string, errs ...error) fiber.Map {
	var styles []CssLink
	styles = append(styles, CssLink{
		Link: "https://getbootstrap.com/docs/5.1/examples/sign-in/signin.css",
	})
	return fiber.Map{
		"CssStyles": styles,
		"Code":      code,
		"Errors":    errs,
	}
}

func ErrorViewData(code int, message string) fiber.Map {
	var styles []CssLink
	return fiber.Map{
		"CssStyles": styles,
		"Code":      code,
		"Message":   message,
	}
}
