package handlers

import (
	"strings"

	"github.com/MiG-21/go-sso/internal"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

const (
	prefix       = "Bearer "
	ctxUserIdKey = "__ctx__user__id__key__"
)

func CtxClaims(ctx *fiber.Ctx) *internal.SignInClaims {
	l := ctx.Locals(ctxUserIdKey)
	if l != nil {
		return l.(*internal.SignInClaims)
	}
	return nil
}

func Authenticate(config *internal.Config) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tokenString := ctx.Get("Authorization")
		if tokenString == "" {
			return fiber.NewError(fiber.StatusBadRequest, "token required")
		}
		if strings.HasPrefix(tokenString, prefix) {
			tokenString = strings.TrimPrefix(tokenString, prefix)
		}
		parsedToken, err := jwt.ParseWithClaims(tokenString, &internal.SignInClaims{}, func(token *jwt.Token) (interface{}, error) {
			// since we only use the one private key to sign the tokens,
			// we also only use its public counterpart to verify
			return config.Crypto.PublicKey, nil
		})
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}
		if claims, ok := parsedToken.Claims.(*internal.SignInClaims); ok && parsedToken.Valid {
			ctx.Locals(ctxUserIdKey, claims)
		} else {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid token")
		}
		return ctx.Next()
	}
}
