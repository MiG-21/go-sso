package sso

type (
	// UserModel data item
	UserModel struct {
		Id       int64  `db:"id,primarykey,autoincrement"`
		Name     string `db:"name,size:255"`
		Email    string `db:"email,primarykey,size:255"`
		Password string `db:"password,size:100"`
		Gender   string `db:"gender,size:50"`
		Data     string `db:"data,size:2048"`
		Created  int64  `db:"created"`
		Updated  int64  `db:"updated"`
		Active   bool   `db:"active"`
		Locked   bool   `db:"locked"`
		LockedTo int64  `db:"locked_to"`
	}

	UserManager interface {
		Authenticate(string, string) (*UserModel, error)
		Create(*UserModel) error
		Validate(*UserModel) error
		ById(int64) (*UserModel, error)
	}
)
