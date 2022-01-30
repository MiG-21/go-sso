package models

type (
	ApplicationModel struct {
		Id          int64  `db:"id,primarykey,autoincrement"`
		Application string `db:"application,size:255"`
		Domain      string `db:"domain,size:255"`
		RedirectUrl string `db:"redirect_url,size:255"`
		Code        string `db:"code,size:50"`
		Created     int64  `db:"created_at"`
		Updated     int64  `db:"updated_at"`
	}

	ApplicationManager interface {
		Create(*ApplicationModel) error
		Update(*ApplicationModel) (int64, error)
		Delete(*ApplicationModel) (int64, error)
		ById(int64) (*ApplicationModel, error)
		ByCode(string) (*ApplicationModel, error)
		List() ([]*ApplicationModel, error)
	}
)
