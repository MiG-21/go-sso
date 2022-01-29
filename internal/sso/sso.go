package sso

import (
	"time"

	"github.com/MiG-21/go-sso/internal"
	"github.com/gofiber/fiber/v2"
)

type (
	// SSOer is what it needs to be implemented for sso functionality.
	SSOer interface {
		UserManager() UserManager
		// CTValidHours returns the cookie/jwt token validity in hours.
		CTValidHours() int64
		CookieName() string
		CookieDomain() string
		// BuildJWTToken takes the user and the user roles info which is then signed by the private
		// key of the login server. The expiry of the token is set per the third argument.
		BuildJWTToken(int64, []string, time.Time) (string, error)
		// BuildCookie takes the jwt token and returns a cookie and sets the expiration time of the same to that of
		// the second arg.
		BuildCookie(string, time.Time) *fiber.Cookie
		// Logout sets the expiration time of the cookie in the past rendering it unusable.
		Logout(time.Time) *fiber.Cookie
	}

	SSO struct {
		Crypto *internal.ConfigCrypto
		Cookie *internal.ConfigCookie
	}
)

func (sso SSO) BuildJWTToken(id int64, roles []string, exp time.Time) (string, error) {
	return internal.GenSignInJWT(id, roles, sso.Crypto.PrivateKey, exp.Unix())
}

func (sso SSO) CTValidHours() int64 {
	return sso.Cookie.ValidHours
}

func (sso SSO) CookieName() string {
	return sso.Cookie.Name
}

func (sso SSO) CookieDomain() string {
	return sso.Cookie.Domain
}

func (sso SSO) BuildCookie(s string, exp time.Time) *fiber.Cookie {
	c := &fiber.Cookie{
		Name:     sso.Cookie.Name,
		Value:    s,
		Domain:   sso.Cookie.Domain,
		Path:     "/",
		Expires:  exp,
		MaxAge:   int(sso.Cookie.ValidHours * 3600),
		Secure:   true,
		HTTPOnly: true,
	}
	return c
}

func (sso SSO) Logout(exp time.Time) *fiber.Cookie {
	c := &fiber.Cookie{
		Name:     sso.Cookie.Name,
		Value:    "",
		Domain:   sso.Cookie.Domain,
		Path:     "/",
		Expires:  exp,
		MaxAge:   -1,
		Secure:   true,
		HTTPOnly: true,
	}
	return c
}

func SetupSSO(config *internal.Config) *SSO {
	return &SSO{&config.Crypto, &config.Cookie}
}
