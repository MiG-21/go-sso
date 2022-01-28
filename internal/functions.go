package internal

import (
	"golang.org/x/crypto/bcrypt"
)

func GetPasswordHash(password []byte) string {
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hashedPassword)
}
