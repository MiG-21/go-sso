package internal

import (
	"crypto/rsa"

	"github.com/dgrijalva/jwt-go"
)

type (
	SignInClaims struct {
		Id    int64    `json:"Id"`
		Roles []string `json:"roles,omitempty"`
		jwt.StandardClaims
	}

	VerificationClaims struct {
		Id     string `json:"id"`
		Action string `json:"action"`
		jwt.StandardClaims
	}
)

func (sic *SignInClaims) IsAuthorized(roles ...string) bool {
	if len(roles) == 0 {
		return true
	}
	for _, role := range roles {
		if ok, _ := InArray(role, sic.Roles); ok {
			return true
		}
	}
	return false
}

func GenSignInJWT(id int64, roles []string, p *rsa.PrivateKey, t int64) (string, error) {
	claims := SignInClaims{
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

func GenVerificationJWT(id, action string, p *rsa.PrivateKey, t int64) (string, error) {
	claims := VerificationClaims{
		id,
		action,
		jwt.StandardClaims{
			ExpiresAt: t,
			Issuer:    "Login_Server",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)

	return token.SignedString(p)
}
