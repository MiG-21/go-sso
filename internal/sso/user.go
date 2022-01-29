package sso

import (
	"crypto/rsa"
	"net/url"
	"time"

	"github.com/MiG-21/go-sso/internal"
	"github.com/gofiber/fiber/v2"
)

type (
	// UserModel data item
	UserModel struct {
		Id        int64  `db:"id,primarykey,autoincrement"`
		Name      string `db:"name,size:255"`
		Email     string `db:"email,primarykey,size:255"`
		Password  string `db:"password,size:100"`
		Gender    string `db:"gender,size:50"`
		Data      string `db:"data,size:2048"`
		Active    bool   `db:"active"`
		Locked    bool   `db:"locked"`
		LockedTo  int64  `db:"locked_to"`
		Code      string `db:"verification_code,size:50"`
		Created   int64  `db:"created_at"`
		Updated   int64  `db:"updated_at"`
		LastVisit int64  `db:"last_visit_at"`
	}

	UserManager interface {
		Authenticate(string, string) (*UserModel, error)
		Create(*UserModel) error
		Update(*UserModel) error
		Validate(*UserModel) error
		ById(int64) (*UserModel, error)
		Verify(string) (bool, error)
	}
)

func (u UserModel) GetVerificationUrl(ctx *fiber.Ctx, p *rsa.PrivateKey) (*url.URL, error) {
	vUrl := &url.URL{}
	vUrl.Scheme = ctx.Protocol()
	vUrl.Host = ctx.Hostname()
	vUrl.Path = "v1/user/verification"
	exp := time.Now().Add(time.Hour * time.Duration(24)).UTC()
	token, err := internal.GenVerificationJWT(u.Code, p, exp.Unix())
	if err != nil {
		return nil, err
	}
	vUrl.RawQuery = "token=" + token
	return vUrl, nil
}
