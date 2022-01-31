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

	UserPasswordRecover struct {
		UserName        string
		UserEmail       string
		VerificationUrl *url.URL
	}
)

func (uc *UserCreated) Name() string {
	return UserCreatedEvent
}

func (uc *UserCreated) Data() interface{} {
	return uc
}

func (uc *UserPasswordRecover) Name() string {
	return PasswordRecoverEvent
}

func (uc *UserPasswordRecover) Data() interface{} {
	return uc
}
