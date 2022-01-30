package dao

import (
	"database/sql"
	"errors"
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

func (u *UserStore) Validate(user *models.UserModel) error {
	query := "SELECT COUNT(*) FROM `users` WHERE `email`=?"
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
	query := "SELECT * FROM `users` WHERE `email`=? LIMIT 1"
	item := &models.UserModel{}
	err := u.db.SelectOne(item, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
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
	return u.db.Insert(user)
}

func (u *UserStore) Update(user *models.UserModel) (int64, error) {
	return u.db.Update(user)
}

func (u *UserStore) Verify(code string) (bool, error) {
	query := "UPDATE `users` SET `active`=1, `verification_code`='', `updated_at`=? WHERE `verification_code`=? LIMIT 1"
	_, err := u.db.Exec(query, time.Now().Unix(), code)
	switch err {
	case sql.ErrNoRows:
		return false, nil
	case nil:
		return true, nil
	default:
		return false, err
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
