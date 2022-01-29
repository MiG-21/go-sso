package dao

import (
	"time"

	"database/sql"
	"github.com/MiG-21/go-sso/internal"
	"github.com/MiG-21/go-sso/internal/sso"
	"go.uber.org/dig"
)

type (
	SetupResult struct {
		dig.Out

		SSOer sso.SSOer
		Error error `group:"errors"`
	}
)

func SetupMysqlDao(config *internal.Config) SetupResult {
	sr := SetupResult{}
	s := sso.SetupSSO(config)

	db, err := sql.Open("mysql", config.Mysql.Dsn)
	if err != nil {
		sr.Error = err
		return sr
	}
	db.SetMaxOpenConns(config.Mysql.MaxOpenConns)
	db.SetMaxIdleConns(config.Mysql.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(config.Mysql.MaxLifetime) * time.Second)

	uStore, err := setupUserStore(db, 10)
	if err != nil {
		sr.Error = err
		return sr
	}

	sr.SSOer = &MysqlDao{
		SSO:       s,
		UserStore: uStore,
	}

	return sr
}
