package main

import (
	"log"

	"github.com/MiG-21/go-sso/internal"
	"github.com/MiG-21/go-sso/internal/dao"
	"github.com/MiG-21/go-sso/internal/event"
	"github.com/MiG-21/go-sso/internal/mail"
	"github.com/MiG-21/go-sso/internal/web"
	"go.uber.org/dig"
)

// @title Swagger go-sso
// @version develop
// @description go-sso
// @BasePath /v1
func main() {
	c := dig.New()

	wrapError(c.Provide(internal.SetupValidator))
	wrapError(c.Provide(internal.SetupConfig))
	wrapError(c.Provide(event.SetupEventService))
	wrapError(c.Provide(internal.SetupLogger))
	wrapError(c.Provide(dao.SetupMysqlDao))
	wrapError(c.Provide(web.SetupServer))
	wrapError(c.Provide(mail.SetupService))

	if err := c.Invoke(internal.Bootstrap); err != nil {
		log.Fatal(err)
	}
}

func wrapError(e error) {
	if e != nil {
		panic(e)
	}
}
