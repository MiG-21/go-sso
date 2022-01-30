package dao

import (
	"io"

	"github.com/MiG-21/go-sso/internal/models"
	"github.com/go-gorp/gorp/v3"
	_ "github.com/go-sql-driver/mysql"
)

type (
	Store struct {
		tableName string
		db        *gorp.DbMap
		stdout    io.Writer
	}

	MysqlDao struct {
		*models.SSO
		UserStore        *UserStore
		ApplicationStore *ApplicationStore
	}
)

func (sso MysqlDao) ApplicationManager() models.ApplicationManager {
	return sso.ApplicationStore
}

func (sso MysqlDao) UserManager() models.UserManager {
	return sso.UserStore
}
