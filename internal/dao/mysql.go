package dao

import (
	"io"

	"database/sql"
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

func (s *Store) execute(query string, args ...interface{}) (sql.Result, error) {
	ret, err := s.db.Exec(query, args)
	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		return ret, nil
	default:
		return nil, err
	}
}

func (sso MysqlDao) ApplicationManager() models.ApplicationManager {
	return sso.ApplicationStore
}

func (sso MysqlDao) UserManager() models.UserManager {
	return sso.UserStore
}
