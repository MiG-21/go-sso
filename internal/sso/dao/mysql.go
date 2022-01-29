package dao

import (
	"database/sql"
	"errors"
	"io"
	"os"
	"time"

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

func (s *UserStore) Update(user *sso.UserModel) error {
	query := "UPDATE `users` SET `name`=?, `gender`=?, `data`=?, `updated_at`=? WHERE id=? LIMIT 1"
	_, err := s.db.Exec(query, user.Name, user.Gender, user.Data, user.Updated, user.Id)
	if err != nil && err == sql.ErrNoRows {
		return nil
	}
	return err
}

func (s *UserStore) Verify(code string) (bool, error) {
	query := "UPDATE `users` SET `active`=1, `verification_code`='', `updated_at`=? WHERE `verification_code`=? LIMIT 1"
	_, err := s.db.Exec(query, time.Now().Unix(), code)
	switch err {
	case sql.ErrNoRows:
		return false, nil
	case nil:
		return true, nil
	default:
		return false, err
	}
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
	table.AddIndex("idx_verification", "Btree", []string{"verification_code"})

	if err := store.db.CreateTablesIfNotExists(); err != nil {
		return nil, err
	}

	_ = store.db.CreateIndex()

	return store, nil
}
