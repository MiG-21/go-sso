package dao

import (
	"errors"
	"io"
	"os"
	"time"

	"database/sql"
	"github.com/MiG-21/go-sso/internal"
	"github.com/MiG-21/go-sso/internal/sso"
	"github.com/go-gorp/gorp/v3"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
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

func (sso MysqlDao) UserManager() sso.UserManager {
	return sso.UserStore
}

func (s *UserStore) ById(id int64) (*sso.UserModel, error) {
	query := "SELECT * FROM `users` WHERE `id`=? LIMIT 1"
	item := &sso.UserModel{}
	err := s.db.SelectOne(item, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return item, nil
}

func (s *UserStore) Validate(user *sso.UserModel) error {
	query := "SELECT COUNT(*) FROM `users` WHERE `email`=?"
	n, err := s.db.SelectInt(query, user.Email)
	if err != nil {
		return err
	}
	if n > 0 {
		return errors.New("email already taken")
	}
	return nil
}

func (s *UserStore) Authenticate(u string, p string) (*sso.UserModel, error) {
	query := "SELECT * FROM `users` WHERE `email`=? LIMIT 1"
	item := &sso.UserModel{}
	err := s.db.SelectOne(item, query, u)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(item.Password), []byte(p)); err != nil {
		return nil, errors.New("invalid credentials")
	}
	return item, nil
}

func (s *UserStore) Create(user *sso.UserModel) error {
	return s.db.Insert(user)
}

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

	table := store.db.AddTableWithName(sso.UserModel{}, store.tableName)
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
