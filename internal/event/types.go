package event

import (
	"net/url"
)

type (
	UserCreated struct {
		UserName        string
		UserEmail       string
		VerificationUrl *url.URL
	}
)

func (uc *UserCreated) Name() string {
	return "user_created"
}

func (uc *UserCreated) Data() interface{} {
	return uc
}
