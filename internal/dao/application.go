package dao

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/MiG-21/go-sso/internal/models"
	"github.com/go-gorp/gorp/v3"
)

type (
	ApplicationStore struct {
		Store
	}
)

func (a *ApplicationStore) Delete(model *models.ApplicationModel) (int64, error) {
	return a.db.Delete(model)
}

func (a *ApplicationStore) Create(model *models.ApplicationModel) error {
	return a.db.Insert(model)
}

func (a *ApplicationStore) Update(model *models.ApplicationModel) (int64, error) {
	return a.db.Update(model)
}

func (a *ApplicationStore) ByCode(code string) (*models.ApplicationModel, error) {
	query := fmt.Sprintf("SELECT * FROM `%s` WHERE `code`=?", a.tableName)
	app := &models.ApplicationModel{}
	err := a.db.SelectOne(app, query, code)
	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		return app, nil
	default:
		return nil, err
	}
}

func (a *ApplicationStore) ById(id int64) (*models.ApplicationModel, error) {
	item, err := a.db.Get(models.ApplicationModel{}, id)
	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		return item.(*models.ApplicationModel), nil
	default:
		return nil, err
	}
}

func (a *ApplicationStore) List() ([]*models.ApplicationModel, error) {
	// TODO implement me
	panic("implement me")
}

func setupApplicationStore(db *sql.DB) (*ApplicationStore, error) {
	store := &ApplicationStore{
		Store{
			db:        &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{Encoding: "UTF8", Engine: "INNODB"}},
			tableName: "applications",
			stdout:    os.Stderr,
		},
	}

	table := store.db.AddTableWithName(models.ApplicationModel{}, store.tableName).SetKeys(true, "Id")
	table.AddIndex("idx_code", "Btree", []string{"code"}).SetUnique(true)

	if err := store.db.CreateTablesIfNotExists(); err != nil {
		return nil, err
	}

	_ = store.db.CreateIndex()

	return store, nil
}
