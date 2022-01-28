package sso

import "time"

type (
	// StoreUser data item
	StoreUser struct {
		Id       int64     `db:"id,primarykey,autoincrement"`
		Name     string    `db:"name,size:255"`
		Email    string    `db:"email,primarykey,size:255"`
		Password string    `db:"password,size:255"`
		Gender   string    `db:"gender,size:50"`
		Data     string    `db:"data,size:2048"`
		Created  time.Time `db:"created"`
		Updated  time.Time `db:"updated"`
		Active   bool      `db:"active"`
		Locked   bool      `db:"locked"`
		LockedTo time.Time `db:"locked_to"`
	}
)
