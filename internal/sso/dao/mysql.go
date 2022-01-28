package dao

import (
	"io"
	"os"
	"time"

	"database/sql"
	"github.com/MiG-21/go-sso/internal"
	"github.com/MiG-21/go-sso/internal/sso"
	"github.com/go-gorp/gorp/v3"
	_ "github.com/go-sql-driver/mysql"
)

type (
	Store struct {
		tableName string
		db        *gorp.DbMap
		stdout    io.Writer
		ticker    *time.Ticker
	}

	UserStore struct {
		Store
	}

	MysqlDao struct {
		*sso.SSO
		UserStore *UserStore
	}
)

func setupUserStore(db *sql.DB, gcInterval int) (*UserStore, error) {
	store := &UserStore{
		Store{
			db:        &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{Encoding: "UTF8", Engine: "INNODB"}},
			tableName: "users",
			stdout:    os.Stderr,
		},
	}
	interval := 600
	if gcInterval > 0 {
		interval = gcInterval
	}
	store.ticker = time.NewTicker(time.Second * time.Duration(interval))

	table := store.db.AddTableWithName(sso.StoreUser{}, store.tableName)
	table.AddIndex("idx_login", "Btree", []string{"email", "password"})

	if err := store.db.CreateTablesIfNotExists(); err != nil {
		return nil, err
	}

	_ = store.db.CreateIndex()

	return store, nil
}

func SetupMysqlDao(config *internal.Config) (sso.SSOer, error) {
	s, err := sso.SetupSSO(config)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("mysql", config.Mysql.Dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(config.Mysql.MaxOpenConns)
	db.SetMaxIdleConns(config.Mysql.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(config.Mysql.MaxLifetime) * time.Second)

	uStore, err := setupUserStore(db, 10)
	if err != nil {
		return nil, err
	}

	return &MysqlDao{
		SSO:       s,
		UserStore: uStore,
	}, nil
}

func (sso MysqlDao) Login(u string, p string) (*sso.StoreUser, error) {
	return sso.UserStore.Get(u, p)
}

func (sso MysqlDao) UserInfo() {

}

func (s *UserStore) Get(u string, p string) (*sso.StoreUser, error) {
	query := "SELECT * FROM `users` WHERE `email`=? AND `password`=? LIMIT 1"
	var item *sso.StoreUser
	err := s.db.SelectOne(item, query, u, p)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return item, nil
}
