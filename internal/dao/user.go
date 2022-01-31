package dao

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/MiG-21/go-sso/internal/models"
	"github.com/go-gorp/gorp/v3"
	"golang.org/x/crypto/bcrypt"
)

type (
	UserStore struct {
		Store
	}
)

func (u *UserStore) ById(id int64) (*models.UserModel, error) {
	item, err := u.db.Get(models.UserModel{}, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return item.(*models.UserModel), nil
}

func (u *UserStore) ByEmail(email string) (*models.UserModel, error) {
	query := fmt.Sprintf("SELECT * FROM `%s` WHERE `email`=? LIMIT 1", u.tableName)
	return u.selectOne(query, email)
}

func (u *UserStore) ByCode(code string) (*models.UserModel, error) {
	query := fmt.Sprintf("SELECT * FROM `%s` WHERE `verification_code`=? LIMIT 1", u.tableName)
	return u.selectOne(query, code)
}

func (u *UserStore) Validate(user *models.UserModel) error {
	query := fmt.Sprintf("SELECT COUNT(*) FROM `%s` WHERE `email`=?", u.tableName)
	n, err := u.db.SelectInt(query, user.Email)
	if err != nil {
		return err
	}
	if n > 0 {
		return errors.New("email already taken")
	}
	return nil
}

func (u *UserStore) Authenticate(email string, password string) (*models.UserModel, error) {
	item, err := u.ByEmail(email)
	if err = bcrypt.CompareHashAndPassword([]byte(item.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}
	if !item.Active {
		return nil, errors.New("user not verified")
	}
	if item.Locked {
		return nil, errors.New("user is locked")
	}
	now := time.Now()
	if item.LockedTo > now.Unix() {
		return nil, errors.New("user is locked to " + now.Format(time.RFC822))
	}
	return item, nil
}

func (u *UserStore) Create(user *models.UserModel) error {
	user.Created = time.Now().Unix()
	return u.db.Insert(user)
}

func (u *UserStore) Update(user *models.UserModel) (int64, error) {
	user.Updated = time.Now().Unix()
	return u.db.Update(user)
}

func (u *UserStore) Delete(model *models.UserModel) error {
	// TODO implement me
	panic("implement me")
}

func (u *UserStore) selectOne(query string, args ...interface{}) (*models.UserModel, error) {
	item := &models.UserModel{}
	err := u.db.SelectOne(item, query, args...)
	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		return item, nil
	default:
		return nil, err
	}
}

func setupUserStore(db *sql.DB) (*UserStore, error) {
	store := &UserStore{
		Store{
			db:        &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{Encoding: "UTF8", Engine: "INNODB"}},
			tableName: "users",
			stdout:    os.Stderr,
		},
	}

	table := store.db.AddTableWithName(models.UserModel{}, store.tableName).SetKeys(true, "Id")
	table.AddIndex("idx_email", "Btree", []string{"email"}).SetUnique(true)
	table.AddIndex("idx_login", "Btree", []string{"email", "password"})
	table.AddIndex("idx_verification", "Btree", []string{"verification_code"})

	if err := store.db.CreateTablesIfNotExists(); err != nil {
		return nil, err
	}

	_ = store.db.CreateIndex()

	return store, nil
}
