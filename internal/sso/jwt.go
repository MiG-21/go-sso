package sso

import (
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"
)

type (
	CustomClaims struct {
		Id    int64    `json:"id"`
		Roles []string `json:"roles,omitempty"`
		jwt.StandardClaims
	}
)

// GenJWT generates the jwt token. Among other stuff, it packs in the authenticated username and the roles that the
// user belongs to and an expiration time. The info is then signed by the private key of the login server.
func GenJWT(id int64, roles []string, p *rsa.PrivateKey, t int64) (string, error) {
	claims := CustomClaims{
		id,
		roles,
		jwt.StandardClaims{
			ExpiresAt: t,
			Issuer:    "Login_Server",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)

	return token.SignedString(p)
}
